package milestone

import (
	"context"
	"fmt"
	"net/http"

	"github.com/417-72KI/gh-milestone/milestone/internal/api"
	"github.com/417-72KI/gh-milestone/milestone/internal/browser"
	"github.com/417-72KI/gh-milestone/milestone/internal/ghrepo"
	iMilestone "github.com/417-72KI/gh-milestone/milestone/internal/milestone"

	"github.com/MakeNowJust/heredoc/v2"

	prShared "github.com/cli/cli/v2/pkg/cmd/pr/shared"
	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/cli/cli/v2/pkg/iostreams"
	"github.com/spf13/cobra"
)

type createOptions struct {
	HttpClient func() (*http.Client, error)
	IO         *iostreams.IOStreams
	Browser    browser.Browser
	Prompter   prShared.Prompt

	Repo ghrepo.Interface

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
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		Browser:    f.Browser,
		Prompter:   f.Prompter,
	}

	createCmd := &cobra.Command{
		Use:   "create [flags]",
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
			opts.Repo = baseRepo

			opts.TitleProvided = cmd.Flags().Changed("title")
			opts.DescriptionProvided = cmd.Flags().Changed("description")
			opts.DueOnProvided = cmd.Flags().Changed("due-on")

			if opts.WebMode && (opts.TitleProvided || opts.DescriptionProvided || opts.DueOnProvided) {
				return cmdutil.FlagErrorf("the `--web` flag is not supported with `--title`, `--description`, or `--due-on`")
			}
			err = createRun(opts)
			if cmdutil.IsUserCancellation(err) {
				cmd.SilenceErrors = true
			}
			return err
		},
	}
	fl := createCmd.Flags()
	fl.StringVarP(&opts.Title, "title", "t", "", "Title for the milestone")
	fl.StringVarP(&opts.Description, "description", "d", "", "Description for the milestone")
	fl.StringVarP(&opts.DueOn, "due-on", "o", "", "Due date for the milestone (format: YYYY/MM/DD)")

	fl.BoolVarP(&opts.WebMode, "web", "w", false, "Open the web browser to create a milestone")
	return createCmd
}

func createRun(opts *createOptions) error {
	if opts.WebMode {
		milestonesURL := ghrepo.GenerateRepoURL(opts.Repo, "milestones/new")
		if opts.IO.IsStdoutTTY() {
			fmt.Fprintf(opts.IO.ErrOut, "Opening %s in your browser.\n", milestonesURL)
		}
		return opts.Browser.Browse(milestonesURL)
	}

	ctx := context.Background()
	opts.IO.DetectTerminalTheme()

	state, err := newMilestoneState(opts)
	if err != nil {
		return err
	}
	message := "\nCreating milestone in %s\n\n"
	cs := opts.IO.ColorScheme()

	if opts.IO.CanPrompt() {
		fmt.Fprintf(opts.IO.ErrOut, message,
			cs.Bold(fmt.Sprintf("%s/%s", opts.Repo.RepoOwner(), opts.Repo.RepoName())))
	}

	if !opts.TitleProvided {
		err = iMilestone.TitleSurvey(opts.Prompter, state)
		if err != nil {
			return err
		}
	}

	if !opts.DescriptionProvided {
		templateContent := ""

		err = iMilestone.DescriptionSurvey(opts.Prompter, state, templateContent)
		if err != nil {
			return err
		}
	}

	if !opts.DueOnProvided {
		err = iMilestone.DueOnSurvey(opts.Prompter, state)
		if err != nil {
			return err
		}
	}

	action, err := iMilestone.ConfirmSubmission(opts.Prompter, false, false)
	if err != nil {
		return fmt.Errorf("unable to confirm: %w", err)
	}
	if action == iMilestone.CancelAction {
		fmt.Fprintln(opts.IO.ErrOut, "Discarding.")
		err = cmdutil.CancelError
		return err
	}

	opts.IO.StartProgressIndicator()
	milestone, err := api.CreateMilestone(ctx, api.CreateMilestoneOptions{
		IO:    opts.IO,
		Repo:  opts.Repo,
		State: state,
	})
	opts.IO.StopProgressIndicator()

	if milestone != nil {
		fmt.Println(*milestone.HTMLURL)
	}
	return err
}

func newMilestoneState(opts *createOptions) (*iMilestone.MilestoneMetadataState, error) {
	state := iMilestone.MilestoneMetadataState{}
	if opts.TitleProvided {
		state.Title = opts.Title
	}
	if opts.DescriptionProvided {
		state.Description = opts.Description
	}
	if opts.DueOnProvided {
		dueOn, err := iMilestone.ParseTime(opts.DueOn)
		if err != nil {
			return nil, err
		}
		state.DueOn = dueOn
	}
	return &state, nil
}
