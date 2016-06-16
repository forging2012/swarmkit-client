package api

import (
	"strings"

	"github.com/docker/swarmkit/api"
)

func parseContainer(cspec *createSpec, spec *api.ServiceSpec) error {
	if len(strings.TrimSpace(cspec.Image)) > 0 {
		spec.Task.GetContainer().Image = cspec.Image
	}

	if len(cspec.Args) > 0 {
		spec.Task.GetContainer().Args = cspec.Args
	}

	if len(cspec.Env) > 0 {
		spec.Task.GetContainer().Env = cspec.Env
	}

	return nil
}
