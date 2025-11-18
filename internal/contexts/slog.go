package contexts

import (
	"context"
	"log/slog"
)

// Logger returns the *slog.Logger from the context.
//
// If the logger is not found, it returns the slog.Default() logger.
//
// The logger is stored in the context with the key (*slog.Logger)(nil).
func Logger(ctx context.Context) *slog.Logger {
	l, ok := ctx.Value((*slog.Logger)(nil)).(*slog.Logger)
	if !ok {
		return slog.Default()
	}

	return l
}

// WithLogger returns a new context with the *slog.Logger.
//
// The logger is stored in the context with the key (*slog.Logger)(nil).
func WithLogger(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, (*slog.Logger)(nil), l)
}
