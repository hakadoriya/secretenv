package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/hakadoriya/secretenv/internal/entrypoint/secretenv"
	"github.com/hakadoriya/secretenv/internal/infra/executor"
	"github.com/hakadoriya/z.go/logz/slogz"
)

func main() {
	os.Exit(Main())
}

func Main() int {
	ctx, stop := signal.NotifyContext(context.Background(), executor.Signals...)
	defer stop()

	l := slog.New(slog.NewTextHandler(os.Stderr, nil))
	ctx = slogz.WithContext(ctx, l)

	if err := secretenv.Entrypoint(ctx, os.Args); err != nil {
		l.Error("secretenv failed", slogz.Error(err))
		return 1
	}

	return 0
}
