package cli

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/rofleksey/serega/internal/app"
	"github.com/rofleksey/serega/internal/config"
	"github.com/spf13/cobra"
)

func newServeCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Serve the HTTP API and frontend",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}

			runtime, err := app.Open(command.Context(), cfg, app.NewLogger(cfg.LogFormat, os.Stderr))
			if err != nil {
				return err
			}

			serveContext, stop := signal.NotifyContext(command.Context(), os.Interrupt, syscall.SIGTERM)
			defer stop()

			return runtime.Serve(serveContext)
		},
	}
}
