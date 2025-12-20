package stdout

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/hakadoriya/secretenv/internal/infra"
	"github.com/hakadoriya/secretenv/internal/infra/internal"
)

// DefaultVersion is an empty string because stdout does not have a version concept
const DefaultVersion = ""

type (
	client struct {
		shell        string
		shellOptions []string
	}
	Option interface {
		apply(c *client)
	}
	optionFunc func(c *client)
)

func (f optionFunc) apply(c *client) { f(c) }

func WithShell(shell string) Option {
	return optionFunc(func(c *client) {
		c.shell = shell
	})
}

func WithShellOptions(shellOptions ...string) Option {
	return optionFunc(func(c *client) {
		c.shellOptions = shellOptions
	})
}

func New(ctx context.Context, opts ...Option) (infra.Client, error) {
	c := &client{
		shell:        "sh",
		shellOptions: []string{"-c"},
	}
	for _, opt := range opts {
		opt.apply(c)
	}
	return c, nil
}

func (c *client) GetSecretStringValue(ctx context.Context, key string, opts ...infra.GetSecretStringValueOption) (value string, err error) {
	// NOTE: stdout does not have a version concept, so cfg is not used,
	// but the argument is received for interface compatibility with other providers.
	cfg := &internal.GetSecretStringValueConfig{
		Version: DefaultVersion,
	}
	for _, opt := range opts {
		opt.Apply(cfg)
	}

	shellArgs := append([]string{c.shell}, c.shellOptions...)
	shellArgs = append(shellArgs, key)

	//nolint:gosec // shell command input is main concept of this provider
	cmd := exec.CommandContext(ctx, shellArgs[0], shellArgs[1:]...)
	cmd.Stdin = os.Stdin
	stdoutBuf := bytes.NewBuffer(nil)
	cmd.Stdout = io.MultiWriter(os.Stdout, stdoutBuf)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("cmd.Run: %w", err)
	}

	return stdoutBuf.String(), nil
}
