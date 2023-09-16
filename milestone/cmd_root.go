package milestone

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewRootCmd(f *cmdutil.Factory) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "milestone",
		Short: "Manage milestones.",
		Long:  `Work with GitHub milestones.`,
		Example: heredoc.Doc(`
			$ gh milestone list
			$ gh milestone view 1
			$ gh milestone create --title "v1.0" --description "Version 1.0"
			$ gh milestone close 1
		`),
	}
	rootCmd.SilenceUsage = true

	cmdutil.EnableRepoOverride(rootCmd, f)

	cmdutil.AddGroup(rootCmd, "General Commands",
		newCreateCmd(f),
		newListCmd(f),
	)
	cmdutil.AddGroup(rootCmd, "Targeted Commands",
		newViewCmd(f),
		newCloseCmd(f),
		newReopenCmd(f),
	)

	return rootCmd
}
