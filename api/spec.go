package api

type (
	createSpec struct {
		Name               string            `json:"name"`                           // service name
		Image              string            `json:"image"`                          // container image
		Labels             map[string]string `json:"labels,omitempty"`               // service label (key=value)
		Mode               string            `json:"mode,omitempty"`                 // one of replicated, global
		Replicas           uint64            `json:"replicas,omitempty"`             // number of replicas for the service (only works in replicated service mode)
		Args               []string          `json:"args,omitempty"`                 // container args
		Env                []string          `json:"env,omitempty"`                  // container env
		Ports              []string          `json:"ports,omitempty"`                // ports
		Network            string            `json:"network,omitempty"`              // network name
		MemoryReservation  string            `json:"memory-reservation,omitempty"`   // amount of reserved memory (e.g. 512m)
		MemoryLimit        string            `json:"memory-limit,omitempty"`         // memory limit (e.g. 512m)
		CPUReservation     string            `json:"cpu-reservation,omitempty"`      // number of CPU cores reserved (e.g. 0.5)
		CPULimit           string            `json:"cpu-limit,omitempty"`            // CPU cores limit (e.g. 0.5)
		UpdateParallelism  uint64            `json:"update-parallelism,omitempty"`   // task update parallelism (0 = all at once)
		UpdateDelay        string            `json:"update-delay,omitempty"`         // delay between task updates (0s = none)
		RestartCondition   string            `json:"restart-condition,omitempty"`    // condition to restart the task (any, failure, none)
		RestartDelay       string            `json:"restart-delay,omitempty"`        // delay between task restarts
		RestartMaxAttempts uint64            `json:"restart-max-attempts,omitempty"` // maximum number of restart attempts (0 = unlimited)
		RestartWindow      string            `json:"restart-window,omitempty"`       // time window to evaluate restart attempts (0 = unbound)
		Constraint         []string          `json:"constraint,omitempty"`           // Placement constraint (node.labels.key==value)
		Bind               []string          `json:"bind,omitempty"`                 // define a bind mount
		Volume             []string          `json:"volume,omitempty"`               // define a volume mount
	}
)
