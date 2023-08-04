package milestone

import (
	"context"
	"fmt"
	"net/http"

	"github.com/417-72KI/gh-milestone/milestone/internal/api"
	"github.com/417-72KI/gh-milestone/milestone/internal/browser"
	"github.com/417-72KI/gh-milestone/milestone/internal/ghrepo"
	iMilestone "github.com/417-72KI/gh-milestone/milestone/internal/milestone"

	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/cli/cli/v2/pkg/iostreams"
	"github.com/cli/go-gh/v2"

	"github.com/spf13/cobra"
)

type issuesOptions struct {
	HttpClient func() (*http.Client, error)
	IO         *iostreams.IOStreams

	BaseRepo ghrepo.Interface

	Exporter cmdutil.Exporter
	Browser  browser.Browser

	Selector string
	WebMode  bool
}

func newIssuesCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &issuesOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		Browser:    f.Browser,
	}

	cmd := &cobra.Command{
		Use:   "issues <number> [flags]",
		Short: "List issues of milestone",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Selector = args[0]
			baseRepo, err := f.BaseRepo()
			if err != nil {
				return err
			}
			opts.BaseRepo = baseRepo

			return issuesRun(opts)
		},
	}
	return cmd
}

func issuesRun(opts *issuesOptions) error {
	ctx := context.Background()
	opts.IO.DetectTerminalTheme()

	opts.IO.StartProgressIndicator()
	num, repo, err := iMilestone.MilestoneNumberAndRepoFromArg(opts.Selector)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}
	if repo != nil {
		opts.BaseRepo = repo
	}

	opts.IO.StartProgressIndicator()
	milestone, err := api.GetMilestone(ctx, opts.BaseRepo, num)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	flags := []string{
		"-R", fmt.Sprintf("%s/%s", opts.BaseRepo.RepoOwner(), opts.BaseRepo.RepoName()),
	}

	issueCommand := append(
		[]string{
			"issue",
			"list",
			"--state",
			"all",
			"--milestone",
			fmt.Sprint(*milestone.Title),
		},
		flags...,
	)
	err = gh.ExecInteractive(ctx, issueCommand...)
	if err != nil {
		return err
	}

	prCommand := append(
		[]string{
			"pr",
			"list",
			"--state",
			"all",
			"-S",
			fmt.Sprintf("milestone:%s", *milestone.Title),
		},
		flags...,
	)

	err = gh.ExecInteractive(ctx, prCommand...)
	if err != nil {
		return err
	}
	return nil
}
