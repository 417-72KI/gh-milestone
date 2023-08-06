package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/417-72KI/gh-milestone/milestone"
	"github.com/cli/cli/v2/pkg/cmd/factory"
	"github.com/cli/go-gh/v2"
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
	cmdFactory := factory.New(version)
	rootCmd, err := milestone.NewRootCmd(cmdFactory)
	if err != nil {
		return exitStatusError
	}
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
	viewOut := strings.Split(stdOut.String(), "\n")[0]
	semver := semverRE.FindStringSubmatch(viewOut)[0]

	return semver, nil
}
