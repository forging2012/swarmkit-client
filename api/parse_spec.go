package api

import (
	"fmt"
	"strings"

	"github.com/docker/swarmkit/api"
)

func merge(cspec *createSpec, spec *api.ServiceSpec, c api.ControlClient) (err error) {
	if len(strings.TrimSpace(cspec.Name)) > 0 {
		spec.Annotations.Name = cspec.Name
	}
	if len(cspec.Labels) > 0 {
		spec.Annotations.Labels = cspec.Labels
	}
	if err = parseMode(cspec, spec); err != nil {
		return
	}
	if err = parseContainer(cspec, spec); err != nil {
		return
	}
	if err = parsePorts(cspec, spec); err != nil {
		return
	}
	if err := parseNetworks(cspec, spec, c); err != nil {
		return err
	}
	if err = parseRestart(cspec, spec); err != nil {
		return err
	}
	if err := parseUpdate(cspec, spec); err != nil {
		return err
	}
	if err := parsePlacement(cspec, spec); err != nil {
		return err
	}
	if err := parseBind(cspec, spec); err != nil {
		return err
	}

	if err := parseVolume(cspec, spec); err != nil {
		return err
	}

	return
}

func parseMode(cspec *createSpec, spec *api.ServiceSpec) error {
	if len(strings.TrimSpace(cspec.Mode)) > 0 {
		switch cspec.Mode {
		case "global":
			if spec.GetGlobal() == nil {
				spec.Mode = &api.ServiceSpec_Global{
					Global: &api.GlobalService{},
				}
			}
		case "replicated":
			if spec.GetReplicated() == nil {
				spec.Mode = &api.ServiceSpec_Replicated{
					Replicated: &api.ReplicatedService{},
				}
			}
		}
	}

	if cspec.Replicas > 0 {
		if spec.GetReplicated() == nil {
			return fmt.Errorf("--replicas can only be specified in --mode replicated")
		}
		spec.GetReplicated().Replicas = cspec.Replicas
	}

	return nil
}
