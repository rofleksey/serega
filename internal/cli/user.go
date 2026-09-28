package cli

import "github.com/spf13/cobra"

func newUserCommand() *cobra.Command {
	command := &cobra.Command{Use: "user", Short: "Manage shared-board accounts"}
	command.AddCommand(newUserCreateCommand())

	return command
}
