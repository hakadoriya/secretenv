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

const DefaultVersion = "AWSCURRENT"

type client struct {
	client *secretsmanager.Client
}

var _ infra.Client = (*client)(nil)

func New(ctx context.Context) (infra.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("config.LoadDefaultConfig: %w", err)
	}

	c := secretsmanager.NewFromConfig(cfg)

	return &client{client: c}, nil
}

func (c *client) GetSecretStringValue(ctx context.Context, key string, opts ...infra.GetSecretStringValueOption) (value string, err error) {
	cfg := &internal.GetSecretStringValueConfig{
		Version: DefaultVersion,
	}
	for _, opt := range opts {
		opt.Apply(cfg)
	}

	//nolint:exhaustruct
	out, err := c.client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(key),
		VersionStage: aws.String(cfg.Version),
	})
	if err != nil {
		return "", fmt.Errorf("c.client.GetSecretValue: %w", err)
	}

	if out.SecretString == nil {
		return "", nil
	}

	return *out.SecretString, nil
}
