package secretenv

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/hakadoriya/z.go/cliz"
	"golang.org/x/sys/unix"

	"github.com/hakadoriya/secretenv/internal/dotenv"
	"github.com/hakadoriya/secretenv/internal/infra"
	"github.com/hakadoriya/secretenv/internal/infra/aws/secretsmanager"
	"github.com/hakadoriya/secretenv/pkg/errors"
)

const (
	optProvider      = "provider"
	envProvider      = "SECRETENV_PROVIDER"
	optSecret        = "secret"
	envSecret        = "SECRETENV_SECRET"
	optSecretVersion = "secret-version"
	envSecretVersion = "SECRETENV_SECRET_VERSION"
)

// Entrypoint is the entrypoint for the secretenv command.
//
// It parses the command line arguments and executes the command.
func Entrypoint(ctx context.Context, args []string) error {
	//nolint:exhaustruct
	cmd := cliz.Command{
		Name:        "secretenv",
		Description: "A tool to manage secrets",
		Options: []cliz.Option{
			//nolint:exhaustruct
			&cliz.StringOption{
				Name:        optProvider,
				Env:         envProvider,
				Description: "The provider to use",
				Required:    true,
			},
			//nolint:exhaustruct
			&cliz.StringOption{
				Name:        optSecret,
				Env:         envSecret,
				Description: "The secret name contains the .env file",
				Required:    true,
			},
			//nolint:exhaustruct
			&cliz.StringOption{
				Name:        optSecretVersion,
				Env:         envSecretVersion,
				Description: "The secret version to use",
			},
		},
		ExecFunc: func(cmd *cliz.Command, args []string) error {
			provider, err := cmd.GetOptionString(optProvider)
			if err != nil {
				return fmt.Errorf("cmd.GetOptionString: %w", err)
			}
			secret, err := cmd.GetOptionString(optSecret)
			if err != nil {
				return fmt.Errorf("cmd.GetOptionString: %w", err)
			}

			var opts []infra.GetSecretStringValueOption
			secretVersion, err := cmd.GetOptionString(optSecretVersion)
			if err != nil {
				return fmt.Errorf("cmd.GetOptionString: %w", err)
			}
			if secretVersion != "" {
				opts = append(opts, infra.WithGetSecretStringValueOptionVersion(secretVersion))
			}

			// 1st argument is the command name, so skip it
			args = args[1:]
			if len(args) < 1 {
				return errors.ErrNoArguments
			}

			var secretClient infra.Client
			switch provider {
			case "aws":
				secretClient, err = secretsmanager.New(ctx)
				if err != nil {
					return fmt.Errorf("secretsmanager.New: %w", err)
				}
			default:
				return fmt.Errorf("provider=%s: %w", provider, errors.ErrUnknownProvider)
			}

			secretValue, err := secretClient.GetSecretStringValue(ctx, secret, opts...)
			if err != nil {
				return fmt.Errorf("secretClient.GetSecretStringValue: %w", err)
			}

			parser, err := dotenv.NewParser(ctx)
			if err != nil {
				return fmt.Errorf("dotenv.NewParser: %w", err)
			}

			dotenv, err := parser.Parse(ctx, secretValue)
			if err != nil {
				return fmt.Errorf("parser.Parse: %w", err)
			}

			envs := os.Environ()
			for _, env := range dotenv.Env {
				envs = append(envs, fmt.Sprintf("%s=%s", env.Key, env.Value))
			}

			execPath, err := exec.LookPath(args[0])
			if err != nil {
				return fmt.Errorf("exec.LookPath: %w", err)
			}

			if err := unix.Exec(execPath, args, envs); err != nil {
				return fmt.Errorf("unix.Exec: %w", err)
			}

			return nil
		},
	}

	if err := cmd.Exec(ctx, args); err != nil {
		return fmt.Errorf("cmd.Exec: %w", err)
	}

	return nil
}
