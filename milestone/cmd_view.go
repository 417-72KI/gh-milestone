package milestone

import (
	"context"
	"fmt"
	"strconv"

	"github.com/417-72KI/gh-milestone/milestone/internal/api"
	iMilestone "github.com/417-72KI/gh-milestone/milestone/internal/milestone"

	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func newViewCmd(f *cmdutil.Factory) *cobra.Command {
	var (
		selector string
		webMode  bool
	)

	io := f.IOStreams

	viewCmd := &cobra.Command{
		Use:   "view <number> [flags]",
		Short: "Display the information about a milestone",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			selector = args[0]

			ctx := context.Background()
			io.DetectTerminalTheme()
			if num, err := strconv.Atoi(selector); err == nil {
				baseRepo, err := f.BaseRepo()
				if err != nil {
					return err
				}
				owner := baseRepo.RepoOwner()
				repo := baseRepo.RepoName()
				io.StartProgressIndicator()
				milestone, err := api.GetMilestone(ctx, owner, repo, num)
				io.StopProgressIndicator()
				if err != nil {
					return err
				}
				if webMode {
					milestoneURL := *milestone.HTMLURL
					if err != nil {
						return err
					}
					if f.IOStreams.IsStdoutTTY() {
						fmt.Fprintf(f.IOStreams.ErrOut, "Opening %s in your browser.\n", milestoneURL)
					}
					f.Browser.Browse(milestoneURL)
					return nil
				}
				if io.IsStdoutTTY() {
					return iMilestone.PrintReadableMilestonePreview(io, milestone)
				}
				return iMilestone.PrintRawMilestonePreview(io.Out, milestone)
			} else {
				return err
			}
		},
	}
	viewCmd.Flags().BoolVarP(&webMode, "web", "w", false, "List milestones in the web browser")
	return viewCmd
}
