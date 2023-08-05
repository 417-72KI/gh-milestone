package milestone

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/417-72KI/gh-milestone/milestone/internal/api"
	iMilestone "github.com/417-72KI/gh-milestone/milestone/internal/milestone"

	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/cli/cli/v2/pkg/iostreams"
	"github.com/spf13/cobra"
)

type viewOptions struct {
	HttpClient func() (*http.Client, error)
	IO         *iostreams.IOStreams

	Selector string
	WebMode  bool
	Exporter cmdutil.Exporter
}

func newViewCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &viewOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	viewCmd := &cobra.Command{
		Use:   "view <number> [flags]",
		Short: "Display the information about a milestone",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Selector = args[0]

			ctx := context.Background()
			opts.IO.DetectTerminalTheme()
			if num, err := strconv.Atoi(opts.Selector); err == nil {
				baseRepo, err := f.BaseRepo()
				if err != nil {
					return err
				}
				owner := baseRepo.RepoOwner()
				repo := baseRepo.RepoName()
				opts.IO.StartProgressIndicator()
				milestone, err := api.GetMilestone(ctx, owner, repo, num)
				opts.IO.StopProgressIndicator()
				if err != nil {
					return err
				}
				if opts.WebMode {
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
				if opts.Exporter != nil {
					output := api.ConvertMilestoneToMap(milestone, opts.Exporter.Fields())
					return opts.Exporter.Write(opts.IO, output)
				}
				if opts.IO.IsStdoutTTY() {
					return iMilestone.PrintReadableMilestonePreview(opts.IO, milestone)
				}
				return iMilestone.PrintRawMilestonePreview(opts.IO.Out, milestone)
			} else {
				return err
			}
		},
	}
	viewCmd.Flags().BoolVarP(&opts.WebMode, "web", "w", false, "List milestones in the web browser")
	cmdutil.AddJSONFlags(viewCmd, &opts.Exporter, api.MilestoneFields)
	return viewCmd
}
