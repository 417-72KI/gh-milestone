package milestone

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/417-72KI/gh-milestone/milestone/internal/api"

	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func newReopenCmd(f *cmdutil.Factory) *cobra.Command {
	var (
		selector string
	)
	cs := f.IOStreams.ColorScheme()

	reopenCmd := &cobra.Command{
		Use:   "reopen {<number> | <url>}",
		Short: "Reopen milestone",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			selector = args[0]

			ctx := context.Background()
			if num, err := strconv.Atoi(selector); err == nil {
				baseRepo, err := f.BaseRepo()
				if err != nil {
					return err
				}
				owner := baseRepo.RepoOwner()
				repo := baseRepo.RepoName()
				f.IOStreams.DetectTerminalTheme()

				f.IOStreams.StartProgressIndicator()
				milestone, err := api.GetMilestone(ctx, owner, repo, num)
				f.IOStreams.StopProgressIndicator()
				if err != nil {
					return err
				}

				f.IOStreams.StartProgressIndicator()
				result, err := api.ReopenMilestone(ctx, api.ReopenMilestoneOptions{
					IO:        f.IOStreams,
					Owner:     owner,
					Repo:      repo,
					Milestone: milestone,
				})
				f.IOStreams.StopProgressIndicator()
				if err != nil {
					return err
				}
				if result != nil {
					fmt.Printf(cs.Green("%s reopend."), *result.HTMLURL)
				}
			} else if url, err := url.Parse(selector); err == nil {
				milestone, err := api.GetMilestoneByURL(ctx, url)
				if err != nil {
					return err
				}
				owner, repo, err := api.FetchOwnerAndRepoFromURL(url)
				if err != nil {
					return err
				}

				result, err := api.ReopenMilestone(ctx, api.ReopenMilestoneOptions{
					IO:        f.IOStreams,
					Owner:     *owner,
					Repo:      *repo,
					Milestone: milestone,
				})
				if err != nil {
					return err
				}
				if result != nil {
					fmt.Printf(cs.Green("%s reopend."), *result.HTMLURL)
				}
			} else {
				return err
			}

			return nil
		},
	}

	return reopenCmd
}
