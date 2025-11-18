package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/hakadoriya/secretenv/internal/contexts"
	"github.com/hakadoriya/secretenv/internal/entrypoint/secretenv"
)

func main() {
	os.Exit(Main(context.Background()))
}

func Main(ctx context.Context) int {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	l := slog.New(slog.NewTextHandler(os.Stderr, nil))
	ctx = contexts.WithLogger(ctx, l)

	if err := secretenv.Entrypoint(ctx, os.Args); err != nil {
		l.Error("secretenv failed", "error", err)
		return 1
	}

	return 0
}
