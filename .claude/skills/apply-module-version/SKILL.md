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
`OLD=github.com/foo/bar` / `NEW=github.com/foo/bar/v2`. This matters for the boundary
rule in step 2: when `$OLD` has no suffix, it is a literal string prefix of `$NEW`.

### 2. Replace the imports

List the affected files:

```sh
rg -l --glob '*.go' "$OLD"
```

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

Match `$OLD` only where the next character is `/` or `"` — that alone stops a shorter
major from eating a longer one (`.../v2` would otherwise corrupt `.../v20/pkg`).

If `$OLD` has no version suffix (the v0/v1 case above), add one more condition: the `/`
must **not** be followed by `v` + digits. Without it, `$OLD` matches its own already-migrated
form — `github.com/foo/bar` is a literal prefix of `github.com/foo/bar/v2/subpkg` — and
rewriting that occurrence corrupts it into `github.com/foo/bar/v2/v2/subpkg`. Verify matches
before editing:

```sh
rg -Pl "\Q$OLD\E(?!/v[0-9])(/|\")" --glob '*.go'
```

Use the Edit tool (`sed -i` needs `sed -i ''` on macOS and is incompatible with GNU sed).

### 3. Run go mod tidy to drop the old major

**Always run this after the replacement.** go.sum doesn't have an `h1:` hash for the new
major yet (only a `/go.mod` line) — building or testing before the replacement will fail.

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

### 5. Commit

Skip committing if there's nothing to commit. Split into two commits to match past convention.

```sh
git add -u -- '*.go'
git commit -m 'replace import'

git add go.mod go.sum
git commit -m 'run `go mod tidy`'
```

`git add -u -- '*.go'` stages every modified tracked `.go` file, so any call-site fixes
from step 4 are picked up automatically alongside the import replacements.

## Definition of done

- No old major left in `go.mod`'s require block
- No `.go` file imports the old major
- `go build ./...` and `go vet ./...` both pass
- The binary built via `make build` runs `./gh-milestone list` successfully
- `git status --short` shows no uncommitted changes

```sh
go mod edit -json | jq -e --arg old "$OLD" '.Require[] | select(.Path == $old)' \
  && echo 'FAIL: old major still in go.mod' || echo 'OK: no old major in go.mod'
rg -Pl "\Q$OLD\E(?!/v[0-9])(/|\")" --glob '*.go' \
  && echo 'FAIL: old major still imported' || echo 'OK: no old major in imports'
go build ./... && go vet ./... && echo 'OK: build/vet passed'
make build && ./gh-milestone list && echo 'OK: runtime check passed'
git status --short
```

Plain `rg "$OLD" go.mod` / `rg "$OLD" --glob '*.go'` would false-positive whenever `$OLD`
has no version suffix, since `$OLD` is then a literal prefix of `$NEW` and matches every
correctly-migrated line too. The `jq` check matches the exact `go.mod` require path instead
of a substring, and the `.go` check reuses the same negative-lookahead rule from step 2.
