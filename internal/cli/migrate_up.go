package cli

import (
	"os"

	"github.com/rofleksey/serega/internal/app"
	"github.com/spf13/cobra"
)

func newMigrateUpCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "up",
		Short: "Apply all pending migrations",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			return app.MigrateUp(command.Context(), os.Getenv("DATABASE_URL"))
		},
	}
}
