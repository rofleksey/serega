package cli

import (
	"errors"

	"github.com/spf13/cobra"
)

func newRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:           "serega",
		Short:         "Serega shared Kanban board",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(*cobra.Command, []string) error {
			return errors.New("a command is required")
		},
	}
	root.AddCommand(newServeCommand(), newMigrateCommand(), newUserCommand())

	return root
}
