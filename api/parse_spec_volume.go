package api

import (
	"fmt"
	"strings"

	"github.com/docker/swarmkit/api"
)

// parseVolume only supports a very simple version of annonymous volumes for
// testing the most basic of data flows. Replace with a --mount flag, similar
// to what we have in docker service.
func parseVolume(cspec *createSpec, spec *api.ServiceSpec) error {
	if len(cspec.Volume) > 0 {
		container := spec.Task.GetContainer()

		for _, volume := range cspec.Volume {
			if strings.Contains(volume, ":") {
				return fmt.Errorf("volume format %q not supported", volume)
			}
			container.Mounts = append(container.Mounts, api.Mount{
				Type:     api.MountTypeVolume,
				Target:   volume,
				Writable: true,
			})
		}
	}

	return nil
}
