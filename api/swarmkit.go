package api

import (
	"crypto/tls"
	"net"
	"time"

	"github.com/docker/swarmkit/api"
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
