package milestone

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/417-72KI/gh-milestone/milestone/internal/api"
	"github.com/417-72KI/gh-milestone/milestone/internal/ghrepo"
	iMilestone "github.com/417-72KI/gh-milestone/milestone/internal/milestone"

	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/cli/cli/v2/pkg/iostreams"
	"github.com/spf13/cobra"
)

type viewOptions struct {
	HttpClient func() (*http.Client, error)
	IO         *iostreams.IOStreams

	Exporter      cmdutil.Exporter
	OpenInBrowser func(string) error

	Selector string
	WebMode  bool
}

func newViewCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &viewOptions{
		IO:            f.IOStreams,
		HttpClient:    f.HttpClient,
		OpenInBrowser: f.Browser.Browse,
	}

	viewCmd := &cobra.Command{
		Use:   "view <number> [flags]",
		Short: "Display the information about a milestone",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Selector = args[0]

			baseRepo, err := f.BaseRepo()
			if err != nil {
				return err
			}
			return viewRun(baseRepo, opts)
		},
	}
	viewCmd.Flags().BoolVarP(&opts.WebMode, "web", "w", false, "List milestones in the web browser")
	cmdutil.AddJSONFlags(viewCmd, &opts.Exporter, api.MilestoneFields)
	return viewCmd
}

func viewRun(repo ghrepo.Interface, opts *viewOptions) error {
	ctx := context.Background()
	opts.IO.DetectTerminalTheme()
	if num, err := strconv.Atoi(opts.Selector); err == nil {
		opts.IO.StartProgressIndicator()
		milestone, err := api.GetMilestone(ctx, repo, num)
		opts.IO.StopProgressIndicator()
		if err != nil {
			return err
		}
		if opts.WebMode {
			milestoneURL := *milestone.HTMLURL
			if err != nil {
				return err
			}
			if opts.IO.IsStdoutTTY() {
				fmt.Fprintf(opts.IO.ErrOut, "Opening %s in your browser.\n", milestoneURL)
			}
			opts.OpenInBrowser(milestoneURL)
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
}
