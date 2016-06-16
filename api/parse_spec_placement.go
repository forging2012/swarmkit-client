package api

import "github.com/docker/swarmkit/api"

func parsePlacement(cspec *createSpec, spec *api.ServiceSpec) error {
	if len(cspec.Constraint) > 0 {
		if spec.Task.Placement == nil {
			spec.Task.Placement = &api.Placement{}
		}
		spec.Task.Placement.Constraints = cspec.Constraint
	}

	return nil
}
