package main

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/417-72KI/gh-milestone/milestone"
	ghcontext "github.com/cli/cli/v2/context"
	"github.com/cli/cli/v2/git"
	"github.com/cli/cli/v2/pkg/cmd/factory"
	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/cli/cli/v2/pkg/iostreams"
	"github.com/cli/go-gh/v2"
	"github.com/cli/go-gh/v2/pkg/ssh"
)

type exitCode int

const (
	// exitStatusOK is status code zero
	exitStatusOK exitCode = iota
	// exitStatusError is status code non-zero
	exitStatusError
)

func main() {
	os.Exit(int(run()))
}

func run() exitCode {
	version, err := ghVersion()
	if err != nil {
		return exitStatusError
	}

	// Construct Factory manually using public APIs (factory.New requires internal types).
	io := iostreams.System()
	gitClient := &git.Client{
		GhPath: "gh",
		Stderr: io.ErrOut,
		Stdin:  io.In,
		Stdout: io.Out,
	}

	// Build Remotes function.
	remotesFn := func() (ghcontext.Remotes, error) {
		ctx := context.Background()
		gitRemotes, err := gitClient.Remotes(ctx)
		if err != nil {
			return nil, err
		}
		if len(gitRemotes) == 0 {
			return nil, fmt.Errorf("no git remotes found")
		}
		translator := ssh.NewTranslator()
		remotes := ghcontext.TranslateRemotes(gitRemotes, translator)
		sort.Sort(remotes)
		return remotes, nil
	}

	// Construct Factory.
	f := &cmdutil.Factory{
		IOStreams:      io,
		AppVersion:     version,
		InvokingAgent:  "",
		ExecutablePath: "gh",
	}
	f.Remotes = remotesFn
	f.BaseRepo = factory.BaseRepoFunc(f.Remotes)

	rootCmd := milestone.NewRootCmd(f)
	if err := rootCmd.Execute(); err != nil {
		return exitStatusError
	}
	return exitStatusOK
}

var semverRE = regexp.MustCompile(`\d+\.\d+\.\d+`)

func ghVersion() (string, error) {
	args := []string{"version"}
	stdOut, _, err := gh.Exec(args...)
	if err != nil {
		return "", fmt.Errorf("failed to view repo: %w", err)
	}
	viewOut := strings.TrimSuffix(stdOut.String(), "\n")
	semver := semverRE.FindStringSubmatch(viewOut)[0]

	return semver, nil
}
