package selector

import (
	"errors"
	"math/rand"
	"sync"

	"github.com/cyub/hyper/pkg/registry"
)

// RoundRobin strategy
func RoundRobin(services *registry.Service) Next {
	nodes := make([]*registry.Node, 0, len(services.Nodes))
	for _, node := range services.Nodes {
		nodes = append(nodes, node)
	}

	var mtx sync.RWMutex
	var i = rand.Int()
	return func() (*registry.Node, error) {
		if len(nodes) == 0 {
			return nil, errors.New("None avaliable node")
		}
		mtx.RLock()
		node := nodes[i%len(nodes)]
		i++
		mtx.RUnlock()
		return node, nil
	}
}

// Strategy type
type Strategy func(services *registry.Service) Next
