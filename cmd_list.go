package milestones

import (
	"context"

	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	var (
		closed     bool
		repository string
		webMode    bool
	)
	listCmd := &cobra.Command{
		Use:   "list [flags]",
		Short: "List milestones in a repository",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			owner, repo, err := parseRepository(repository)
			if err != nil {
				return err
			}

			ctx := context.Background()
			milestones, err := milestones(ctx, owner, repo, closed)
			cmd.Print(milestones)

			return nil
		},
	}

	listCmd.Flags().BoolVarP(&closed, "closed", "c", false, "Show closed milestones.")
	listCmd.Flags().BoolVarP(&webMode, "web", "w", false, "List milestones in the web browser")
	listCmd.Flags().StringVarP(&repository, "repo", "R", "", "Select another repository using the OWNER/REPO format")

	return listCmd
}

func matchAll(checks ...cobra.PositionalArgs) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		for _, check := range checks {
			if err := check(cmd, args); err != nil {
				return err
			}
		}
		return nil
	}
}
