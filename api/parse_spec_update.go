package api

import (
	"strings"
	"time"

	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/protobuf/ptypes"
)

func parseUpdate(cspec *createSpec, spec *api.ServiceSpec) error {
	if cspec.UpdateParallelism > 0 {
		if spec.Update == nil {
			spec.Update = &api.UpdateConfig{}
		}
		spec.Update.Parallelism = cspec.UpdateParallelism
	}

	if len(strings.TrimSpace(cspec.UpdateDelay)) > 0 {
		delayDuration, err := time.ParseDuration(cspec.UpdateDelay)
		if err != nil {
			return err
		}

		if spec.Update == nil {
			spec.Update = &api.UpdateConfig{}
		}
		spec.Update.Delay = *ptypes.DurationProto(delayDuration)
	}
	return nil
}
