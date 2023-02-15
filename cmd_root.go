package milestones

import (
	"github.com/spf13/cobra"
)

func NewRootCmd() (*cobra.Command, error) {
	rootCmd := &cobra.Command {
		Use: "milestones",
		Short: "Create, Edit, List milestones.",
	}
	rootCmd.AddCommand(newListCmd())
	
	return rootCmd, nil
}