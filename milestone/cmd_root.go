package milestone

import (
	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewRootCmd(f *cmdutil.Factory) (*cobra.Command, error) {
	rootCmd := &cobra.Command{
		Use:   "milestones",
		Short: "Manage milestones.",
		Long:  `Work with GitHub milestones.`,
	}
	cmdutil.EnableRepoOverride(rootCmd, f)

	cmdutil.AddGroup(rootCmd, "General commands",
		newListCmd(f),
		newViewCmd(f),
		newCreateCmd(f),
		newCloseCmd(f),
		newReopenCmd(f),
	)

	return rootCmd, nil
}
