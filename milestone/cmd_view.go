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
	"github.com/spf13/cobra"
)

type viewOptions struct {
	HttpClient func() (*http.Client, error)
	IO         *iostreams.IOStreams
	BaseRepo   ghrepo.Interface

	Exporter cmdutil.Exporter
	Browser  browser.Browser

	Selector string
	WebMode  bool
}

func newViewCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &viewOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		Browser:    f.Browser,
	}

	viewCmd := &cobra.Command{
		Use:   "view {<number> | <url>} [flags]",
		Short: "Display the information about a milestone",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Selector = args[0]

			baseRepo, err := f.BaseRepo()
			if err != nil {
				return err
			}
			opts.BaseRepo = baseRepo
			return viewRun(opts)
		},
	}
	viewCmd.Flags().BoolVarP(&opts.WebMode, "web", "w", false, "List milestones in the web browser")
	cmdutil.AddJSONFlags(viewCmd, &opts.Exporter, api.MilestoneFields)
	return viewCmd
}

func viewRun(opts *viewOptions) error {
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
	if opts.WebMode {
		milestoneURL := *milestone.HTMLURL
		if opts.IO.IsStdoutTTY() {
			fmt.Fprintf(opts.IO.ErrOut, "Opening %s in your browser.\n", milestoneURL)
		}
		return opts.Browser.Browse(milestoneURL)
	}
	if opts.Exporter != nil {
		output := api.ConvertMilestoneToMap(milestone, opts.Exporter.Fields())
		return opts.Exporter.Write(opts.IO, output)
	}
	if opts.IO.IsStdoutTTY() {
		return iMilestone.PrintReadableMilestonePreview(opts.IO, milestone)
	}
	return iMilestone.PrintRawMilestonePreview(opts.IO.Out, milestone)
}
