package milestone

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func newCloseCmd(f *cmdutil.Factory) *cobra.Command {
	var (
		selector string
	)
	cs := f.IOStreams.ColorScheme()

	closeCmd := &cobra.Command{
		Use:   "close <number>", /* "close {<number> | <url>}", */
		Short: "Close milestone",
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
				milestone, err := getMilestone(ctx, owner, repo, num)
				if err != nil {
					return err
				}
				result, err := closeMilestone(closeMilestoneOptions{
					ctx:       ctx,
					IO:        f.IOStreams,
					owner:     owner,
					repo:      repo,
					milestone: milestone,
				})
				if err != nil {
					return err
				}
				fmt.Printf(cs.Green("%s closed."), *result.HTMLURL)
			} else if url, err := url.Parse(selector); err == nil {
				return fmt.Errorf("closing by URL not supported yet. %s", url)
				// milestone, err := getMilestoneByURL(ctx, url)
				// if err != nil {
				// 	return err
				// }
				// result, err = closeMilestone(closeMilestoneOptions{
				// 	ctx:       ctx,
				// 	IO:        f.IOStreams,
				// 	milestone: milestone,
				// })
				// if err != nil {
				// 	return err
				// }
				// fmt.Printf(cs.Green("%s closed."), *result.HTMLURL)
			} else {
				return err
			}

			return nil
		},
	}

	return closeCmd
}
