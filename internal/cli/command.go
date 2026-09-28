// Package cli owns Serega's command-line surface.
package cli

import (
	"context"
)

// Execute runs the Serega CLI with explicit arguments.
func Execute(ctx context.Context, args []string) error {
	command := newRootCommand()
	command.SetArgs(args)

	return command.ExecuteContext(ctx)
}
