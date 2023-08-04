package milestone

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/cli/cli/v2/pkg/iostreams"
	"github.com/spf13/cobra"
)

type listOptions struct {
	HttpClient func() (*http.Client, error)
	IO         *iostreams.IOStreams

	State   string
	WebMode bool
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
				milestonesURL := GenerateRepositoryURL(host, owner, repo, "milestones")
				if f.IOStreams.IsStdoutTTY() {
					fmt.Fprintf(f.IOStreams.ErrOut, "Opening %s in your browser.\n", milestonesURL)
				}
				f.Browser.Browse(milestonesURL)
				return nil
			}

			ctx := context.Background()
			f.IOStreams.DetectTerminalTheme()

			f.IOStreams.StartProgressIndicator()
			listResult, err := milestones(ctx, owner, repo, opts.State)
			f.IOStreams.StopProgressIndicator()
			if err != nil {
				return err
			}
			if len(listResult) == 0 {
				switch opts.State {
				case "open":
					fmt.Fprintf(f.IOStreams.Out, "no open milestones in %s/%s", owner, repo)
				default:
					fmt.Fprintf(f.IOStreams.Out, "no milestones match your search in %s/%s", owner, repo)
				}
			} else {
				PrintMilestones(f.IOStreams, time.Now(), "", len(listResult), listResult)
			}

			return nil
		},
	}

	cmdutil.StringEnumFlag(listCmd, &opts.State, "state", "s", "open", []string{"open", "closed", "all"}, "Filter by state")
	listCmd.Flags().BoolVarP(&opts.WebMode, "web", "w", false, "List milestones in the web browser")

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
