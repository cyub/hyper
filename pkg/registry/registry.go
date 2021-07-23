package registry

import "context"

// Registry inferface
type Registry interface {
	Register(context.Context, *Service) error
	Deregister(context.Context, *Service) error
	GetService(context.Context, string) (*Service, error)
	String() string
}

// Service struct
type Service struct {
	Name     string            `json:"name"`
	Nodes    []*Node           `json:"nodes"`
	Metadata map[string]string `json:"metadata"`
}

// Node struct
type Node struct {
	ID       string            `json:"id"`
	Address  string            `json:"address"`
	Metadata map[string]string `json:"metadata"`
}
