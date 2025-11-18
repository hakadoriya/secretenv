package infra

import (
	"context"

	"github.com/hakadoriya/secretenv/internal/infra/internal"
)

type GetSecretStringValueOption interface {
	Apply(c *internal.GetSecretStringValueConfig)
}

type getSecretStringValueOptionFunc func(c *internal.GetSecretStringValueConfig)

func (f getSecretStringValueOptionFunc) Apply(c *internal.GetSecretStringValueConfig) { f(c) }

func WithGetSecretStringValueOptionVersion(version string) GetSecretStringValueOption {
	return getSecretStringValueOptionFunc(func(c *internal.GetSecretStringValueConfig) {
		c.Version = version
	})
}

type Client interface {
	GetSecretStringValue(ctx context.Context, key string, opts ...GetSecretStringValueOption) (value string, err error)
}
