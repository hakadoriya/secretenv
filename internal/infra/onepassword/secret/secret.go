package secret

import (
	"context"
	"fmt"
	"os"

	"github.com/1password/onepassword-sdk-go"

	"github.com/hakadoriya/secretenv/internal/infra"
	"github.com/hakadoriya/secretenv/internal/infra/internal"
	"github.com/hakadoriya/secretenv/pkg/errors"
)

// DefaultVersion is an empty string because 1Password does not have a version concept
const DefaultVersion = ""

// ServiceAccountTokenEnvKey is the environment variable name for the 1Password service account token
//
//nolint:gosec // G101: This is not a credential, but a constant for the environment variable name
const ServiceAccountTokenEnvKey = "OP_SERVICE_ACCOUNT_TOKEN"

// IntegrationName is the integration name for the 1Password SDK
const IntegrationName = "secretenv"

// IntegrationVersion is the integration version for the 1Password SDK
const IntegrationVersion = "v1.0.0"

type (
	client struct {
		serviceAccountToken string
		client              *onepassword.Client
	}
	Option interface {
		apply(c *client)
	}
	optionFunc func(c *client)
)

func (f optionFunc) apply(c *client) { f(c) }

func WithServiceAccountToken(token string) Option {
	return optionFunc(func(c *client) {
		c.serviceAccountToken = token
	})
}

var _ infra.Client = (*client)(nil)

// New initializes and returns a 1Password client.
// It retrieves the service account token from the OP_SERVICE_ACCOUNT_TOKEN environment variable.
func New(ctx context.Context, opts ...Option) (infra.Client, error) {
	c := &client{
		serviceAccountToken: os.Getenv(ServiceAccountTokenEnvKey),
		client:              nil,
	}
	for _, opt := range opts {
		opt.apply(c)
	}

	if c.serviceAccountToken == "" {
		return nil, fmt.Errorf("env=%s: %w", ServiceAccountTokenEnvKey, errors.ErrEnvVarNotSet)
	}

	var err error
	c.client, err = onepassword.NewClient(
		ctx,
		onepassword.WithServiceAccountToken(c.serviceAccountToken),
		onepassword.WithIntegrationInfo(IntegrationName, IntegrationVersion),
	)
	if err != nil {
		return nil, fmt.Errorf("onepassword.NewClient: %w", err)
	}

	return c, nil
}

// GetSecretStringValue is a method that retrieves the value of a secret corresponding to the specified key.
// key is expected to be in the format "op://vault/item/field".
//
// NOTE: 1Password does not have a version concept, so the Version option in opts is ignored.
func (c *client) GetSecretStringValue(ctx context.Context, key string, opts ...infra.GetSecretStringValueOption) (value string, err error) {
	// NOTE: 1Password does not have a version concept, so cfg is not used,
	// but the argument is received for interface compatibility with other providers.
	cfg := &internal.GetSecretStringValueConfig{
		Version: DefaultVersion,
	}
	for _, opt := range opts {
		opt.Apply(cfg)
	}

	secret, err := c.client.Secrets().Resolve(ctx, key)
	if err != nil {
		return "", fmt.Errorf("c.client.Secrets().Resolve: %w", err)
	}

	return secret, nil
}
