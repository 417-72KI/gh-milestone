package milestone

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/417-72KI/gh-milestone/milestone/internal/api"
	"github.com/417-72KI/gh-milestone/milestone/internal/browser"
	"github.com/417-72KI/gh-milestone/milestone/internal/ghrepo"
	"github.com/417-72KI/gh-milestone/milestone/internal/milestone"

	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/cli/cli/v2/pkg/iostreams"
	"github.com/spf13/cobra"
)

type listOptions struct {
	HttpClient func() (*http.Client, error)
	IO         *iostreams.IOStreams
	BaseRepo   ghrepo.Interface
	Browser    browser.Browser

	State    string
	WebMode  bool
	Exporter cmdutil.Exporter
}

func newListCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &listOptions{
		HttpClient: f.HttpClient,
		IO:         f.IOStreams,
		Browser:    f.Browser,
	}

	listCmd := &cobra.Command{
		Use:   "list [flags]",
		Short: "List milestones in a repository",
		Args:  cmdutil.NoArgsQuoteReminder,
		RunE: func(cmd *cobra.Command, args []string) error {
			baseRepo, err := f.BaseRepo()
			if err != nil {
				return err
			}
			opts.BaseRepo = baseRepo
			return listRun(opts)
		},
	}

	cmdutil.StringEnumFlag(listCmd, &opts.State, "state", "s", "open", []string{"open", "closed", "all"}, "Filter by state")
	listCmd.Flags().BoolVarP(&opts.WebMode, "web", "w", false, "List milestones in the web browser")
	cmdutil.AddJSONFlags(listCmd, &opts.Exporter, api.MilestoneFields)

	return listCmd
}

func listRun(opts *listOptions) error {
	if opts.WebMode {
		milestonesURL := ghrepo.GenerateRepoURL(opts.BaseRepo, "milestones")
		if opts.IO.IsStdoutTTY() {
			fmt.Fprintf(opts.IO.ErrOut, "Opening %s in your browser.\n", milestonesURL)
		}
		return opts.Browser.Browse(milestonesURL)
	}

	milestoneState := strings.ToLower(opts.State)

	filterOptions := api.FilterOptions{
		State:  milestoneState,
		Fields: []string{},
	}

	if opts.Exporter != nil {
		filterOptions.Fields = opts.Exporter.Fields()
	}

	ctx := context.Background()
	opts.IO.DetectTerminalTheme()

	opts.IO.StartProgressIndicator()
	listResult, err := api.Milestones(ctx, opts.BaseRepo, filterOptions)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}
	if len(listResult) == 0 && opts.Exporter == nil {
		switch opts.State {
		case "open":
			fmt.Fprintf(opts.IO.Out, "no open milestones in %s/%s", opts.BaseRepo.RepoOwner(), opts.BaseRepo.RepoName())
		default:
			fmt.Fprintf(opts.IO.Out, "no milestones match your search in %s/%s", opts.BaseRepo.RepoOwner(), opts.BaseRepo.RepoName())
		}
		return nil
	}
	if err := opts.IO.StartPager(); err == nil {
		defer opts.IO.StopPager()
	} else {
		fmt.Fprintf(opts.IO.ErrOut, "failed to start pager: %v\n", err)
	}

	if opts.Exporter != nil {
		outputs := []map[string]any{}
		for _, result := range listResult {
			output := api.ConvertMilestoneToMap(result, filterOptions.Fields)
			outputs = append(outputs, output)
		}
		return opts.Exporter.Write(opts.IO, outputs)
	}

	milestone.PrintMilestones(opts.IO, time.Now(), "", len(listResult), listResult)

	return nil
}
