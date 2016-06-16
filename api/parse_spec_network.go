package api

import (
	"strings"

	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/cmd/swarmctl/network"
	ct "golang.org/x/net/context"
)

func parseNetworks(cspec *createSpec, spec *api.ServiceSpec, c api.ControlClient) error {
	if len(strings.TrimSpace(cspec.Network)) > 0 {
		n, err := network.GetNetwork(ct.TODO(), c, cspec.Network)
		if err != nil {
			return err
		}

		spec.Networks = []*api.ServiceSpec_NetworkAttachmentConfig{
			{
				Target: n.ID,
			},
		}
	}
	return nil
}
