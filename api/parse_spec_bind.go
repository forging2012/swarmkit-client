package api

import (
	"fmt"
	"strings"

	"github.com/docker/swarmkit/api"
)

// parseBind only supports a very simple version of bind for testing the most
// basic of data flows. Replace with a --mount flag, similar to what we have in
// docker service.
func parseBind(cspec *createSpec, spec *api.ServiceSpec) error {
	if len(cspec.Bind) > 0 {
		container := spec.Task.GetContainer()

		for _, bind := range cspec.Bind {
			parts := strings.SplitN(bind, ":", 2)
			if len(parts) != 2 {
				return fmt.Errorf("bind format %q not supported", bind)
			}
			container.Mounts = append(container.Mounts, api.Mount{
				Type:     api.MountTypeBind,
				Source:   parts[0],
				Target:   parts[1],
				Writable: true,
			})
		}
	}

	return nil
}
