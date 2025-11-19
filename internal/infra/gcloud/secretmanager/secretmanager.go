package secretmanager

import (
	"context"
	"fmt"
	"strings"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/hakadoriya/secretenv/internal/infra"
	"github.com/hakadoriya/secretenv/internal/infra/internal"
)

const DefaultVersion = "latest"

type client struct {
	client *secretmanager.Client
}

var _ infra.Client = (*client)(nil)

func New(ctx context.Context) (infra.Client, error) {
	c, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("secretmanager.NewClient: %w", err)
	}

	return &client{client: c}, nil
}

func (c *client) GetSecretStringValue(ctx context.Context, key string, opts ...infra.GetSecretStringValueOption) (value string, err error) {
	cfg := &internal.GetSecretStringValueConfig{
		Version: DefaultVersion,
	}
	for _, opt := range opts {
		opt.Apply(cfg)
	}

	if parts := strings.Split(key, "/versions/"); len(parts) > 1 {
		// like: `projects/1234567890/secrets/my-secret/versions/1`
		if cfg.Version == DefaultVersion {
			// do nothing
		} else {
			// If cfg.Version is not `latest`, override it
			key = parts[0] + "/versions/" + cfg.Version
		}
	} else {
		// like: `projects/1234567890/secrets/my-secret`
		key += "/versions/" + cfg.Version
	}

	out, err := c.client.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{
		Name: key,
	})
	if err != nil {
		return "", fmt.Errorf("c.client.AccessSecretVersion: %w", err)
	}

	return string(out.Payload.Data), nil
}
