package secretsmanager

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"

	"github.com/hakadoriya/secretenv/internal/infra"
	"github.com/hakadoriya/secretenv/internal/infra/internal"
)

type client struct {
	svc *secretsmanager.Client
}

var _ infra.Client = (*client)(nil)

func New(ctx context.Context) (infra.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("config.LoadDefaultConfig: %w", err)
	}

	svc := secretsmanager.NewFromConfig(cfg)

	return &client{svc: svc}, nil
}

func (c *client) GetSecretStringValue(ctx context.Context, key string, opts ...infra.GetSecretStringValueOption) (value string, err error) {
	cfg := &internal.GetSecretStringValueConfig{
		Version: "AWSCURRENT",
	}
	for _, opt := range opts {
		opt.Apply(cfg)
	}

	//nolint:exhaustruct
	out, err := c.svc.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(key),
		VersionStage: aws.String(cfg.Version),
	})
	if err != nil {
		return "", fmt.Errorf("secretsmanager.GetSecretValue: %w", err)
	}

	if out.SecretString == nil {
		return "", nil
	}

	return *out.SecretString, nil
}
