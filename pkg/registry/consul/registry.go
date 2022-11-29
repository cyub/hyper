// Copyright 2022 tink <qietingfy@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package consul

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/cyub/hyper/pkg/registry"
	consul "github.com/hashicorp/consul/api"
)

const name = "consul"

var (
	// ErrNoneAvailable none node avaliable
	ErrNoneAvailable = errors.New("Require at leaset one node")
)

type consulRegistry struct {
	opts registry.Options
	sync.Mutex
	client *consul.Client
}

func (cr *consulRegistry) Register(ctx context.Context, s *registry.Service) error {
	if len(s.Nodes) == 0 {
		return ErrNoneAvailable
	}

	node := s.Nodes[0]
	var checks consul.AgentServiceChecks
	if cr.opts.Context != nil {
		interval, ok := cr.opts.Context.Value("consul_tcp_check").(time.Duration)
		if ok {
			checks = append(checks, &consul.AgentServiceCheck{
				TCP:      node.Address,
				Interval: fmt.Sprintf("%v", interval),
			})
		}
	}

	check, ok := ctx.Value("service_check").(*consul.AgentServiceCheck)
	if ok {
		checks = append(checks, check)
	}

	host, pt, err := net.SplitHostPort(node.Address)
	if err != nil {
		return err
	}

	port, _ := strconv.Atoi(pt)
	registration := &consul.AgentServiceRegistration{
		Name:    s.Name,
		ID:      node.ID,
		Address: host,
		Port:    port,
		Meta:    node.Metadata,
		Checks:  checks,
	}

	return cr.client.Agent().ServiceRegister(registration)
}

func (cr *consulRegistry) Deregister(ctx context.Context, s *registry.Service) error {
	if len(s.Nodes) == 0 {
		return ErrNoneAvailable
	}
	node := s.Nodes[0]
	return cr.client.Agent().ServiceDeregister(node.ID)
}

func (cr *consulRegistry) GetService(ctx context.Context, name string) (srv *registry.Service, err error) {
	var resp []*consul.ServiceEntry
	queryOptions, ok := ctx.Value("service_query_options").(*consul.QueryOptions)
	if ok && queryOptions != nil {
		resp, _, err = cr.client.Health().Service(name, "", false, queryOptions)
	} else {
		resp, _, err = cr.client.Health().Service(name, "", false, nil)
	}

	if err != nil {
		return nil, err
	}
	for _, res := range resp {
		if res.Service.Service != name {
			continue
		}
		var skip bool
		for _, check := range res.Checks {
			if check.Status == "critical" {
				skip = true
				break
			}
		}
		if skip {
			continue
		}

		if srv == nil {
			srv = &registry.Service{
				Name: res.Service.Service,
			}
		}
		srv.Nodes = append(srv.Nodes, &registry.Node{
			ID:       res.Service.ID,
			Address:  net.JoinHostPort(res.Service.Address, strconv.Itoa(res.Service.Port)),
			Metadata: res.Service.Meta,
		})
	}

	return srv, nil
}

func (cr *consulRegistry) String() string {
	return name
}

// NewRegistry return registry
func NewRegistry(opts ...registry.Option) registry.Registry {
	options := registry.Options{
		Context: context.Background(),
		Timeout: 200 * time.Millisecond,
	}
	for _, opt := range opts {
		opt(&options)
	}
	cr := &consulRegistry{
		opts: options,
	}

	config := consul.DefaultConfig()
	var addrs []string
	for _, addr := range cr.opts.Address {
		if host, port, err := net.SplitHostPort(addr); err == nil {
			addrs = append(addrs, net.JoinHostPort(host, port))
		}
	}

	config.Address = addrs[0]
	if config.HttpClient == nil {
		config.HttpClient = new(http.Client)
	}
	if cr.opts.Timeout > 0 {
		config.HttpClient.Timeout = cr.opts.Timeout
	}
	client, _ := consul.NewClient(config)
	cr.client = client
	return cr
}
