// Package observability connects Serega's domain vocabulary to shared logging.
package observability

import (
	"context"
	"io"
	"log/slog"

	"github.com/happytoolin/unolog"
	"github.com/rofleksey/meg/logging"
	"github.com/rofleksey/meg/operation"
)

func NewLogger(format string, output io.Writer) *slog.Logger { return logging.New(format, output) }
func NewEventRuntime(logger *slog.Logger) *unolog.Runtime    { return operation.NewRuntime(logger) }
func StartOperation(ctx context.Context, runtime *unolog.Runtime, domain unolog.Domain, name, id string, fields ...any) *unolog.Operation {
	return operation.Start(ctx, runtime, domain, name, id, fields...)
}
func Enrich(ctx context.Context, fields ...any)  { operation.Enrich(ctx, fields...) }
func RecordError(ctx context.Context, err error) { operation.RecordError(ctx, err) }
func CompleteOperation(op *unolog.Operation, outcome string, err error) {
	operation.Complete(op, outcome, err)
}
