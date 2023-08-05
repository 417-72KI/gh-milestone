package milestone

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/417-72KI/gh-milestone/milestone/internal/api"
	"github.com/417-72KI/gh-milestone/milestone/internal/milestone"
	"github.com/417-72KI/gh-milestone/milestone/internal/utils"

	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/cli/cli/v2/pkg/iostreams"
	"github.com/spf13/cobra"
)

type listOptions struct {
	HttpClient func() (*http.Client, error)
	IO         *iostreams.IOStreams

	Assignee string
	State    string
	Author   string
	WebMode  bool
	Exporter cmdutil.Exporter
}

func newListCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &listOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
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
			host := baseRepo.RepoHost()
			owner := baseRepo.RepoOwner()
			repo := baseRepo.RepoName()

			if opts.WebMode {
				milestonesURL := utils.GenerateRepositoryURL(host, owner, repo, "milestones")
				if f.IOStreams.IsStdoutTTY() {
					fmt.Fprintf(f.IOStreams.ErrOut, "Opening %s in your browser.\n", milestonesURL)
				}
				f.Browser.Browse(milestonesURL)
				return nil
			}

			milestoneState := strings.ToLower(opts.State)

			filterOptions := api.FilterOptions{
				State:  milestoneState,
				Author: opts.Author,
				Fields: []string{},
			}

			if opts.Exporter != nil {
				filterOptions.Fields = opts.Exporter.Fields()
			}

			ctx := context.Background()
			f.IOStreams.DetectTerminalTheme()

			f.IOStreams.StartProgressIndicator()
			listResult, err := api.Milestones(ctx, owner, repo, filterOptions)
			f.IOStreams.StopProgressIndicator()
			if err != nil {
				return err
			}
			if len(listResult) == 0 && opts.Exporter == nil {
				switch opts.State {
				case "open":
					fmt.Fprintf(f.IOStreams.Out, "no open milestones in %s/%s", owner, repo)
				default:
					fmt.Fprintf(f.IOStreams.Out, "no milestones match your search in %s/%s", owner, repo)
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
					output := map[string]any{}
					for _, field := range filterOptions.Fields {
						switch field {
						case "url":
							output[field] = *result.URL
						case "id":
							output[field] = *result.ID
						case "number":
							output[field] = *result.Number
						case "state":
							output[field] = *result.State
						case "title":
							output[field] = *result.Title
						case "description":
							output[field] = *result.Description
						case "creator":
							output[field] = *result.Creator.Login
						case "openIssues":
							output[field] = *result.OpenIssues
						case "closedIssues":
							output[field] = *result.ClosedIssues
						case "createdAt":
							output[field] = result.CreatedAt.Format(time.RFC3339)
						case "updatedAt":
							output[field] = result.UpdatedAt.Format(time.RFC3339)
						case "closedAt":
							output[field] = result.ClosedAt.Format(time.RFC3339)
						case "dueOn":
							output[field] = *result.DueOn
						}
					}
					outputs = append(outputs, output)
				}
				return opts.Exporter.Write(opts.IO, outputs)
			}

			milestone.PrintMilestones(f.IOStreams, time.Now(), "", len(listResult), listResult)

			return nil
		},
	}

	cmdutil.StringEnumFlag(listCmd, &opts.State, "state", "s", "open", []string{"open", "closed", "all"}, "Filter by state")
	listCmd.Flags().BoolVarP(&opts.WebMode, "web", "w", false, "List milestones in the web browser")
	cmdutil.AddJSONFlags(listCmd, &opts.Exporter, api.MilestoneFields)

	return listCmd
}

func MatchAll(checks ...cobra.PositionalArgs) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		for _, check := range checks {
			if err := check(cmd, args); err != nil {
				return err
			}
		}
		return nil
	}
}
