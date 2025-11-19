package secretenv

import (
	"context"
	"fmt"
	"os"

	"github.com/hakadoriya/z.go/buildinfoz"
	"github.com/hakadoriya/z.go/cliz"

	"github.com/hakadoriya/secretenv/internal/dotenv"
	"github.com/hakadoriya/secretenv/internal/infra"
	"github.com/hakadoriya/secretenv/internal/infra/aws/secretsmanager"
	"github.com/hakadoriya/secretenv/internal/infra/executor"
	"github.com/hakadoriya/secretenv/internal/infra/gcloud/secretmanager"
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
func Entrypoint(ctx context.Context, osArgs []string) error {
	//nolint:exhaustruct
	c := cliz.Command{
		Name:        "secretenv",
		Description: "A command-line tool that fetches secrets from secret management services (such as Google Cloud Secret Manager, AWS Secrets Manager) and executes commands with those secrets as environment variables",
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
		SubCommands: []*cliz.Command{
			{
				Name:        "version",
				Description: "Print the version of secretenv",
				ExecFunc: func(c *cliz.Command, args []string) error {
					if err := buildinfoz.Fprint(c.Stdout()); err != nil {
						return fmt.Errorf("buildinfoz.Fprint: %w", err)
					}
					return nil
				},
			},
		},
		ExecFunc: execFunc(executor.NewExecutor()),
	}

	if err := c.Exec(ctx, osArgs); err != nil {
		return fmt.Errorf("c.Exec: %w", err)
	}

	return nil
}

func execFunc(e executor.Executor) func(cmd *cliz.Command, args []string) error {
	return func(cmd *cliz.Command, args []string) error {
		ctx := cmd.Context()

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

		// 1st argument is the command name (== secretenv), so skip it
		args = args[1:]
		if len(args) < 1 {
			return errors.ErrNoArguments
		}

		var secretClient infra.Client
		switch provider {
		case "aws":
			secretClient, err = secretsmanager.New(ctx)
			if err != nil {
				return fmt.Errorf("provider=%s: secretsmanager.New: %w", provider, err)
			}
		case "gcloud":
			secretClient, err = secretmanager.New(ctx)
			if err != nil {
				return fmt.Errorf("provider=%s: secretmanager.New: %w", provider, err)
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

		if err := e.Exec(args[0], args, envs); err != nil {
			return fmt.Errorf("unix.Exec: %w", err)
		}

		return nil
	}
}
