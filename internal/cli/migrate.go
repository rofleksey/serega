package cli

import "github.com/spf13/cobra"

func newMigrateCommand() *cobra.Command {
	migrate := &cobra.Command{Use: "migrate", Short: "Manage the PostgreSQL schema"}
	migrate.AddCommand(newMigrateUpCommand())

	return migrate
}
