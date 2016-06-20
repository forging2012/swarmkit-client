package swarmkit

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"github.com/docker/swarmkit/api"
	ct "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Dial establishes a connection and creates a client.
// It infers connection parameters from CLI options.
func Dial(socketAddr string) (api.ControlClient, error) {
	opts := []grpc.DialOption{}
	insecureCreds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})
	opts = append(opts, grpc.WithTransportCredentials(insecureCreds))
	opts = append(opts, grpc.WithDialer(
		func(addr string, timeout time.Duration) (net.Conn, error) {
			return net.DialTimeout("unix", addr, timeout)
		}))
	conn, err := grpc.Dial(socketAddr, opts...)
	if err != nil {
		return nil, err
	}

	client := api.NewControlClient(conn)
	return client, nil
}

// GetService get service with service id or service name in cluster
func GetService(ctx ct.Context, c api.ControlClient, input string) (*api.Service, error) {
	// GetService to match via full ID.
	rg, err := c.GetService(ctx, &api.GetServiceRequest{ServiceID: input})
	if err != nil {
		// If any error (including NotFound), ListServices to match via ID prefix and full name.
		rl, err := c.ListServices(ctx,
			&api.ListServicesRequest{
				Filters: &api.ListServicesRequest_Filters{
					Names: []string{input},
				},
			},
		)
		if err != nil {
			return nil, err
		}

		if len(rl.Services) == 0 {
			return nil, fmt.Errorf("service %s not found", input)
		}

		if l := len(rl.Services); l > 1 {
			return nil, fmt.Errorf("service %s is ambiguous (%d matches found)", input, l)
		}

		return rl.Services[0], nil
	}
	return rg.Service, nil
}

// GetNode get node with name or id from cluster
func GetNode(ctx ct.Context, c api.ControlClient, input string) (*api.Node, error) {
	// GetNode to match via full ID.
	rg, err := c.GetNode(ctx, &api.GetNodeRequest{NodeID: input})
	if err != nil {
		// If any error (including NotFound), ListServices to match via ID prefix and full name.
		rl, err := c.ListNodes(ctx,
			&api.ListNodesRequest{
				Filters: &api.ListNodesRequest_Filters{
					Names: []string{input},
				},
			},
		)
		if err != nil {
			return nil, err
		}

		if len(rl.Nodes) == 0 {
			return nil, fmt.Errorf("node %s not found", input)
		}

		if l := len(rl.Nodes); l > 1 {
			return nil, fmt.Errorf("node %s is ambiguous (%d matches found)", input, l)
		}

		return rl.Nodes[0], nil
	}
	return rg.Node, nil
}

// GetNetwork tries to query for a network as an ID and if it can't be
// found tries to query as a name. If the name query returns exactly
// one entry then it is returned to the caller. Otherwise an error is
// returned.
func GetNetwork(ctx ct.Context, c api.ControlClient, input string) (*api.Network, error) {
	// GetService to match via full ID.
	rg, err := c.GetNetwork(ctx, &api.GetNetworkRequest{NetworkID: input})
	if err != nil {
		// If any error (including NotFound), ListServices to match via ID prefix and full name.
		rl, err := c.ListNetworks(ctx,
			&api.ListNetworksRequest{
				Filters: &api.ListNetworksRequest_Filters{
					Names: []string{input},
				},
			},
		)
		if err != nil {
			return nil, err
		}

		if len(rl.Networks) == 0 {
			return nil, fmt.Errorf("network %s not found", input)
		}

		if l := len(rl.Networks); l > 1 {
			return nil, fmt.Errorf("network %s is ambiguous (%d matches found)", input, l)
		}

		return rl.Networks[0], nil
	}

	return rg.Network, nil
}
