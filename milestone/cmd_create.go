package milestone

import (
	"fmt"
	"net/http"

	"github.com/417-72KI/gh-milestone/milestone/internal/utils"
	"github.com/MakeNowJust/heredoc"
	prShared "github.com/cli/cli/v2/pkg/cmd/pr/shared"
	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/cli/cli/v2/pkg/iostreams"
	"github.com/spf13/cobra"
)

type createOptions struct {
	HttpClient    func() (*http.Client, error)
	IO            *iostreams.IOStreams
	OpenInBrowser func(string) error
	Prompter      prShared.Prompt

	TitleProvided       bool
	DescriptionProvided bool
	DueOnProvided       bool

	WebMode bool

	Title       string
	Description string
	DueOn       string
}

func newCreateCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &createOptions{
		IO:            f.IOStreams,
		HttpClient:    f.HttpClient,
		OpenInBrowser: f.Browser.Browse,
		Prompter:      f.Prompter,
	}

	createCmd := &cobra.Command{
		Use:   "create {<number> | <url>}",
		Short: "Create milestone",
		Example: heredoc.Doc(`
			$ gh milestone create --title "v1.0" --description "Version 1.0"
			$ gh milestone create --title "v1.0" --due-on "2022/01/01"
		`),
		Args: cmdutil.NoArgsQuoteReminder,
		RunE: func(cmd *cobra.Command, args []string) error {
			baseRepo, err := f.BaseRepo()
			if err != nil {
				return err
			}
			host := baseRepo.RepoHost()
			owner := baseRepo.RepoOwner()
			repo := baseRepo.RepoName()

			opts.TitleProvided = cmd.Flags().Changed("title")
			opts.DescriptionProvided = cmd.Flags().Changed("description")
			opts.DueOnProvided = cmd.Flags().Changed("due-on")

			if opts.WebMode && (opts.TitleProvided || opts.DescriptionProvided || opts.DueOnProvided) {
				return fmt.Errorf("the `--web` flag is not supported with `--title`, `--description`, or `--due-on`")
			}

			return createRun(host, owner, repo, opts)
		},
	}
	fl := createCmd.Flags()
	fl.StringVarP(&opts.Title, "title", "t", "", "Title for the milestone")
	fl.StringVarP(&opts.Description, "description", "d", "", "Description for the milestone")
	fl.StringVarP(&opts.DueOn, "due-on", "o", "", "Due date for the milestone (format: YYYY/MM/DD)")

	fl.BoolVarP(&opts.WebMode, "web", "w", false, "Open the web browser to create a milestone")
	return createCmd
}

func createRun(host string, owner string, repo string, opts *createOptions) error {
	if opts.WebMode {
		milestonesURL := utils.GenerateRepositoryURL(host, owner, repo, "milestones/new")
		if opts.IO.IsStdoutTTY() {
			fmt.Fprintf(opts.IO.ErrOut, "Opening %s in your browser.\n", milestonesURL)
		}
		return opts.OpenInBrowser(milestonesURL)
	}

	return fmt.Errorf("not implemented")
}
