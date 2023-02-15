package main

import (
	"os"

	milestones "github.com/417-72KI/gh-milestones"
	"github.com/cli/cli/v2/pkg/cmd/factory"
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
	cmdFactory := factory.New("DEV")
	rootCmd, err := milestones.NewRootCmd(cmdFactory)
	if err != nil {
		return exitStatusError
	}
	if err := rootCmd.Execute(); err != nil {
		return exitStatusError
	}
	return exitStatusOK
}
