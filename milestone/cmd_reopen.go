package milestone

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/417-72KI/gh-milestone/milestone/internal/api"
	"github.com/417-72KI/gh-milestone/milestone/internal/ghrepo"

	"github.com/google/go-github/v53/github"

	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/cli/cli/v2/pkg/iostreams"
	"github.com/spf13/cobra"
)

type reopenOptions struct {
	httpClient func() (*http.Client, error)
	IO         *iostreams.IOStreams

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

			ctx := context.Background()
			if num, err := strconv.Atoi(opts.Selector); err == nil {
				baseRepo, err := f.BaseRepo()
				if err != nil {
					return err
				}
				opts.IO.DetectTerminalTheme()

				opts.IO.StartProgressIndicator()
				milestone, err := api.GetMilestone(ctx, baseRepo, num)
				opts.IO.StopProgressIndicator()
				if err != nil {
					return err
				}
				return reopenMilestone(ctx, opts.IO, baseRepo, milestone)
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
				return reopenMilestone(ctx, opts.IO, ghrepo.NewWithHost(*owner, *repo, url.Hostname()), milestone)
			} else {
				return err
			}
		},
	}

	return reopenCmd
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
