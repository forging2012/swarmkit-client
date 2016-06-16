package api

import (
	"fmt"
	"strings"
	"time"

	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/protobuf/ptypes"
)

func parseRestart(cspec *createSpec, spec *api.ServiceSpec) error {
	if len(strings.TrimSpace(cspec.RestartCondition)) > 0 {

		if spec.Task.Restart == nil {
			spec.Task.Restart = &api.RestartPolicy{}
		}

		switch cspec.RestartCondition {
		case "none":
			spec.Task.Restart.Condition = api.RestartOnNone
		case "failure":
			spec.Task.Restart.Condition = api.RestartOnFailure
		case "any":
			spec.Task.Restart.Condition = api.RestartOnAny
		default:
			return fmt.Errorf("invalid restart condition: %s", cspec.RestartCondition)
		}
	}

	if len(strings.TrimSpace(cspec.RestartDelay)) > 0 {
		delayDuration, err := time.ParseDuration(cspec.RestartDelay)
		if err != nil {
			return err
		}

		if spec.Task.Restart == nil {
			spec.Task.Restart = &api.RestartPolicy{}
		}
		spec.Task.Restart.Delay = ptypes.DurationProto(delayDuration)
	}

	if cspec.RestartMaxAttempts > 0 {
		if spec.Task.Restart == nil {
			spec.Task.Restart = &api.RestartPolicy{}
		}
		spec.Task.Restart.MaxAttempts = cspec.RestartMaxAttempts
	}

	if len(strings.TrimSpace(cspec.RestartWindow)) > 0 {
		windowDelay, err := time.ParseDuration(cspec.RestartWindow)
		if err != nil {
			return err
		}

		if spec.Task.Restart == nil {
			spec.Task.Restart = &api.RestartPolicy{}
		}
		spec.Task.Restart.Window = ptypes.DurationProto(windowDelay)
	}

	return nil
}
