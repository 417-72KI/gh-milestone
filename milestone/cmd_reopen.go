package milestone

import (
	"context"
	"fmt"
	"net/http"

	"github.com/417-72KI/gh-milestone/milestone/internal/api"
	"github.com/417-72KI/gh-milestone/milestone/internal/ghrepo"
	iMilestone "github.com/417-72KI/gh-milestone/milestone/internal/milestone"

	"github.com/google/go-github/v55/github"

	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/cli/cli/v2/pkg/iostreams"
	"github.com/spf13/cobra"
)

type reopenOptions struct {
	httpClient func() (*http.Client, error)
	IO         *iostreams.IOStreams
	BaseRepo   ghrepo.Interface

	Selector string
}

func newReopenCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &reopenOptions{
		httpClient: f.HttpClient,
		IO:         f.IOStreams,
	}

	reopenCmd := &cobra.Command{
		Use:   "reopen {<number> | <url>}",
		Short: "Reopen milestone",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Selector = args[0]

			baseRepo, err := f.BaseRepo()
			if err != nil {
				return err
			}
			opts.BaseRepo = baseRepo

			return reopenRun(opts)
		},
	}

	return reopenCmd
}

func reopenRun(opts *reopenOptions) error {
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

	return reopenMilestone(ctx, opts.IO, opts.BaseRepo, milestone)
}

func reopenMilestone(ctx context.Context, io *iostreams.IOStreams, repo ghrepo.Interface, milestone *github.Milestone) error {
	io.StartProgressIndicator()
	result, err := api.ReopenMilestone(ctx, api.ReopenMilestoneOptions{
		IO:        io,
		Repo:      repo,
		Milestone: milestone,
	})
	io.StopProgressIndicator()
	if err != nil {
		return err
	}
	if result != nil {
		fmt.Printf(io.ColorScheme().Green("%s reopend."), *result.HTMLURL)
	}
	return nil
}
