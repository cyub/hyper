package selector

import (
	"context"
	"errors"

	"github.com/cyub/hyper/pkg/registry"
	"github.com/cyub/hyper/pkg/registry/consul"
)

const name = "registry-selector"

// Selector interface
type Selector interface {
	Select(context.Context, string) (Next, error)
	String() string
}

// Next return registry.Node
type Next func() (*registry.Node, error)

type registrySelector struct {
	opts Options
}

func (s *registrySelector) Select(ctx context.Context, name string) (Next, error) {
	service, err := s.opts.Registry.GetService(context.Background(), name)
	if err != nil {
		return nil, err
	}

	if service == nil {
		return nil, errors.New("None avaliable service node")
	}
	strategy := s.opts.Strategy
	if sg, ok := ctx.Value("select_strategy").(Strategy); ok {
		strategy = sg
	}

	return strategy(service), nil
}

func (s *registrySelector) String() string {
	return name
}

// NewSelector select
func NewSelector(opts ...Option) Selector {
	options := Options{
		Strategy: RoundRobin,
	}

	for _, opt := range opts {
		opt(&options)
	}

	if options.Registry == nil {
		options.Registry = consul.NewRegistry()
	}
	selector := &registrySelector{
		opts: options,
	}
	return selector
}
