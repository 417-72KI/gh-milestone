package milestone

import (
	"fmt"
	"net/http"

	"github.com/MakeNowJust/heredoc"
	prShared "github.com/cli/cli/v2/pkg/cmd/pr/shared"
	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/cli/cli/v2/pkg/iostreams"
	"github.com/spf13/cobra"
)

type createOptions struct {
	HttpClient func() (*http.Client, error)
	IO         *iostreams.IOStreams

	WebMode  bool
	Prompter prShared.Prompt
}

func newCreateCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &createOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
		Prompter:   f.Prompter,
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

			return createRun(opts)
		},
	}

	return createCmd
}

func createRun(opts *createOptions) error {
	return fmt.Errorf("not implemented")
}
