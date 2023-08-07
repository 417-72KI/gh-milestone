package milestone

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/417-72KI/gh-milestone/milestone/internal/api"
	"github.com/google/go-github/v53/github"

	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/cli/cli/v2/pkg/iostreams"
	"github.com/spf13/cobra"
)

type closeOptions struct {
	HttpClient func() (*http.Client, error)
	IO         *iostreams.IOStreams

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

			ctx := context.Background()
			if num, err := strconv.Atoi(opts.Selector); err == nil {
				baseRepo, err := f.BaseRepo()
				if err != nil {
					return err
				}
				owner := baseRepo.RepoOwner()
				repo := baseRepo.RepoName()
				opts.IO.DetectTerminalTheme()

				opts.IO.StartProgressIndicator()
				milestone, err := api.GetMilestone(ctx, owner, repo, num)
				opts.IO.StopProgressIndicator()
				if err != nil {
					return err
				}

				return closeMilestone(ctx, opts.IO, owner, repo, milestone)
			} else if url, err := url.Parse(opts.Selector); err == nil {
				opts.IO.DetectTerminalTheme()

				opts.IO.StartProgressIndicator()
				milestone, err := api.GetMilestoneByURL(ctx, url)
				opts.IO.StopProgressIndicator()
				if err != nil {
					return err
				}
				opts.IO.StartProgressIndicator()
				owner, repo, err := api.FetchOwnerAndRepoFromURL(url)
				opts.IO.StopProgressIndicator()
				if err != nil {
					return err
				}
				return closeMilestone(ctx, opts.IO, *owner, *repo, milestone)
			} else {
				return err
			}
		},
	}

	return closeCmd
}

func closeMilestone(ctx context.Context, io *iostreams.IOStreams, owner string, repo string, milestone *github.Milestone) error {
	io.StartProgressIndicator()
	result, err := api.CloseMilestone(ctx, api.CloseMilestoneOptions{
		IO:        io,
		Owner:     owner,
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
