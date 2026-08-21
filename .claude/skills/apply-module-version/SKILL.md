---
name: apply-module-version
description: Fixes the situation where a Renovate PR bumps a Go module's major version in go.mod but the import statements are still on the old major. Use when go.mod's require block has two majors of the same module coexisting (e.g. .../v85 and .../v90).
---

## When to apply

Only run this when the direct require block in `go.mod` lists two majors of the same module path.
If nothing coexists, do nothing and stop.

## Steps

### 1. Identify the target module and its old/new major

Renovate brings the go.mod change already committed, so `git diff -- go.mod` is always empty. Don't use it.

Find modules with coexisting majors (exclude indirect requires, so normal coexistence like `heredoc` v1/v2 isn't picked up):

```sh
go mod edit -json \
  | jq -r '.Require[] | select(.Indirect != true) | .Path' \
  | sed -E 's|/v[0-9]+$||' | sort | uniq -d
```

The path printed is the target. Check what Renovate actually added with:

```sh
git show HEAD -- go.mod
```

Old major = the one the source code still imports, new major = the one Renovate added. Treat these as variables from here on:

```sh
OLD=github.com/google/go-github/v85   # example
NEW=github.com/google/go-github/v90   # example
```

`$OLD` may have no version suffix at all — Go modules only add `/vN` to the import path
starting at v2 (see https://go.dev/blog/v2-go-modules), so a v0/v1 → v2 bump looks like
`OLD=github.com/foo/bar` / `NEW=github.com/foo/bar/v2`. Step 2's replace rule accounts for this.

### 2. Replace the imports

Replace the module path in every import. Two forms both need handling — some modules are
imported only at the root (`github.com/MakeNowJust/heredoc/v2`), some both ways
(`github.com/cli/go-gh/v2` is imported bare in `cmd/gh-milestone/main.go` and as
`/pkg/text` elsewhere):

| form | before | after |
| --- | --- | --- |
| subpackage | `"$OLD/github"` | `"$NEW/github"` |
| module root | `"$OLD"` | `"$NEW"` |

Handling only the subpackage form leaves the old major imported, so `go mod tidy` keeps it
in `go.mod` — exactly the state this skill exists to remove.

List the affected files with the same rule the replacement must follow: match `$OLD` only
where the next character is `/` or `"`, and where that `/` is not followed by `v` + digits.
Both conditions always apply. This needs ripgrep built with PCRE2 (`rg -P`) — check before
relying on it, and self-heal via this repo's `Brewfile` (which declares `ripgrep`, whose
Homebrew formula bundles PCRE2 at build time) if it's missing:

```sh
rg --pcre2-version >/dev/null 2>&1 || brew bundle install
rg --pcre2-version >/dev/null 2>&1 || echo 'MISSING: rg with PCRE2 support, even after brew bundle install'

rg -Pl "\Q$OLD\E(?!/v[0-9])(/|\")" --glob '*.go'
```

`\Q...\E` quotes `$OLD` so the dots in a module path stay literal — a bare `rg "$OLD"` would
treat them as regex wildcards and match unrelated paths. The character boundary stops a
shorter major from eating a longer one (`.../v2` would otherwise corrupt `.../v20/pkg`). The
negative lookahead matters when `$OLD` has no version suffix (the v0/v1 case above):
`github.com/foo/bar` is then a literal prefix of an already-migrated
`github.com/foo/bar/v2/subpkg`, and without the lookahead the replace would corrupt it into
`github.com/foo/bar/v2/v2/subpkg`. It's a no-op when `$OLD` is versioned.

If the preflight still prints `MISSING` after the retry, stop and report it. `rg`'s default
engine supports neither `\Q...\E` nor lookahead, so there's no drop-in substitute for the
pattern — installing the PCRE2 library alone wouldn't have helped either, since PCRE2 support
is compiled into the `rg` binary itself, not added by a separate library install. Fall back to
inspecting the import lines by hand.

Use the Edit tool (`sed -i` needs `sed -i ''` on macOS and is incompatible with GNU sed).

Commit this right away, before step 3. Doing it now — while the only diff is the mechanical
path swap — keeps this commit pure even if step 4 later needs call-site fixes in the same
files (e.g. `internal/api/github.go` importing the new path AND needing its client-construction
code rewritten for a breaking change): those land as their own commit, not lumped in here.

```sh
git add -u -- '*.go'
git commit -m 'replace import'
```

Skip this (and any commit step below) if there's nothing to commit.

### 3. Run go mod tidy to drop the old major

**Run this after the replacement, and before any build.** Renovate adds the new major to
`go.mod`, but go.sum only gets its `/go.mod` line — not the `h1:` zip hash. That's harmless
while the imports still point at the old major (`go build ./...` passes, since the new major
isn't imported yet). The moment step 2 rewrites the imports, the build breaks until tidy
fills in the hash:

```
missing go.sum entry for module providing package github.com/google/go-github/v90/github
```

`-mod=readonly` is the default from Go 1.16 on, so the build errors out instead of quietly
adding the missing hash.

```sh
go mod tidy
```

### 4. Build it and actually run it

```sh
go build ./...
go vet ./...
```

Then confirm the built binary actually works:

```sh
make build
./gh-milestone list
```

`make build` drops the binary at the repo root (gitignored).
Run the binary **directly** rather than `make list` — `make list` doesn't run
`gh extension install`, so it would exercise the already-installed old `gh milestone`
binary instead of the one you just built. This step hits the real GitHub API, so it
needs `gh` auth and network access. If unavailable, skip it and report that you skipped it.

**A major bump can include breaking changes.** If this step fails, a string replacement wasn't enough:

- Check the target module's release notes / CHANGELOG for breaking changes and fix the call sites
- If the fix requires a design decision or is broad in scope, stop and report it instead of
  silently working around it or hacking something together just to make it compile
- Commit the fix on its own, separate from both the `replace import` commit (step 2) and the
  `go mod tidy` commit (step 5), with a message describing the API change (e.g. `fix breaking
  client construction API in go-github v90`) — not a generic message shared with either of those:

```sh
git add -u -- '*.go'
git commit -m 'fix breaking client construction API in go-github v90'   # example message
```

### 5. Commit the go.mod/go.sum update

```sh
git add go.mod go.sum
git commit -m 'run `go mod tidy`'
```

## Definition of done

- No old major left in `go.mod`'s require block
- No `.go` file imports the old major
- `go build ./...` and `go vet ./...` both pass
- The binary built via `make build` runs `./gh-milestone list` successfully
- `git status --short` shows no uncommitted changes

```sh
rg --pcre2-version >/dev/null 2>&1 || brew bundle install
rg --pcre2-version >/dev/null 2>&1 || echo 'WARNING: rg missing or built without PCRE2, even after brew bundle install — the import check below cannot run'

# Both of these must print nothing. Any output means the migration isn't finished.
go mod edit -json | jq -r --arg old "$OLD" '.Require[] | select(.Path == $old) | .Path'
rg -Pl "\Q$OLD\E(?!/v[0-9])(/|\")" --glob '*.go'

go build ./... && go vet ./...
make build && ./gh-milestone list
git status --short
```

Read the output rather than branching on exit status. `cmd && echo FAIL || echo OK` prints
`OK` whenever the tool itself fails — a missing `jq` (exit 127), a bad filter (exit 3), or
`rg -P` on a build without PCRE2 (exit 2) — which hides a check that never actually ran.

The `jq` check matches the exact `go.mod` require path instead of a plain substring grep,
which would false-positive whenever `$OLD` has no version suffix (it's then a literal prefix
of `$NEW` and matches every correctly-migrated line too). The `.go` check reuses the same
negative-lookahead rule from step 2 for the same reason.
