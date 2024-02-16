package milestone

import (
	"context"
	"fmt"
	"net/http"

	"github.com/417-72KI/gh-milestone/milestone/internal/api"
	"github.com/417-72KI/gh-milestone/milestone/internal/ghrepo"
	iMilestone "github.com/417-72KI/gh-milestone/milestone/internal/milestone"

	"github.com/google/go-github/v59/github"

	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/cli/cli/v2/pkg/iostreams"
	"github.com/spf13/cobra"
)

type closeOptions struct {
	HttpClient func() (*http.Client, error)
	IO         *iostreams.IOStreams
	BaseRepo   ghrepo.Interface

	Selector string
}

func newCloseCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &closeOptions{
		HttpClient: f.HttpClient,
		IO:         f.IOStreams,
	}

	closeCmd := &cobra.Command{
		Use:   "close {<number> | <url>}",
		Short: "Close milestone",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Selector = args[0]

			baseRepo, err := f.BaseRepo()
			if err != nil {
				return err
			}
			opts.BaseRepo = baseRepo

			return closeRun(opts)
		},
	}

	return closeCmd
}

func closeRun(opts *closeOptions) error {
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

	return closeMilestone(ctx, opts.IO, opts.BaseRepo, milestone)
}

func closeMilestone(ctx context.Context, io *iostreams.IOStreams, repo ghrepo.Interface, milestone *github.Milestone) error {
	io.StartProgressIndicator()
	result, err := api.CloseMilestone(ctx, api.CloseMilestoneOptions{
		IO:        io,
		Repo:      repo,
		Milestone: milestone,
	})
	io.StopProgressIndicator()
	if err != nil {
		return err
	}
	if result != nil {
		fmt.Printf(io.ColorScheme().Green("%s closed."), *result.HTMLURL)
	}
	return nil
}
