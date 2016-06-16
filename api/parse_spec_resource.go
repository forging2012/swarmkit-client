package api

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/docker/go-units"
	"github.com/docker/swarmkit/api"
)

func parseResource(cspec *createSpec, spec *api.ServiceSpec) error {
	if len(strings.TrimSpace(cspec.MemoryReservation)) > 1 {
		if spec.Task.Resources == nil {
			spec.Task.Resources = &api.ResourceRequirements{}
		}
		if spec.Task.Resources.Reservations == nil {
			spec.Task.Resources.Reservations = &api.Resources{}
		}
		if err := parseResourceMemory(cspec.MemoryReservation, spec.Task.Resources.Reservations); err != nil {
			return err
		}
	}

	if len(strings.TrimSpace(cspec.MemoryLimit)) > 1 {
		if spec.Task.Resources == nil {
			spec.Task.Resources = &api.ResourceRequirements{}
		}
		if spec.Task.Resources.Limits == nil {
			spec.Task.Resources.Limits = &api.Resources{}
		}
		if err := parseResourceMemory(cspec.MemoryLimit, spec.Task.Resources.Limits); err != nil {
			return err
		}
	}

	if len(strings.TrimSpace(cspec.CPUReservation)) > 0 {
		if spec.Task.Resources == nil {
			spec.Task.Resources = &api.ResourceRequirements{}
		}
		if spec.Task.Resources.Reservations == nil {
			spec.Task.Resources.Reservations = &api.Resources{}
		}
		if err := parseResourceCPU(cspec.CPUReservation, spec.Task.Resources.Reservations); err != nil {
			return err
		}
	}

	if len(strings.TrimSpace(cspec.CPULimit)) > 0 {
		if spec.Task.Resources == nil {
			spec.Task.Resources = &api.ResourceRequirements{}
		}
		if spec.Task.Resources.Limits == nil {
			spec.Task.Resources.Limits = &api.Resources{}
		}
		if err := parseResourceCPU(cspec.CPULimit, spec.Task.Resources.Limits); err != nil {
			return err
		}
	}

	return nil
}

func parseResourceCPU(cpu string, resources *api.Resources) error {

	nanoCPUs, ok := new(big.Rat).SetString(cpu)
	if !ok {
		return fmt.Errorf("invalid cpu: %s", cpu)
	}
	cpuRat := new(big.Rat).Mul(nanoCPUs, big.NewRat(1e9, 1))
	if !cpuRat.IsInt() {
		return fmt.Errorf("CPU value cannot have more than 9 decimal places: %s", cpu)
	}
	resources.NanoCPUs = cpuRat.Num().Int64()
	return nil
}

func parseResourceMemory(memory string, resources *api.Resources) error {
	bytes, err := units.RAMInBytes(memory)
	if err != nil {
		return err
	}

	resources.MemoryBytes = bytes
	return nil
}
