package subsystems

type ResourceConfig struct {
	MemoryLimit string
	CpuShare    string
	CpuSet      string
}

type Subsystem interface {
	// Name get name
	Name() string
	// Set set limit
	Set(path string, res *ResourceConfig) error
	// Apply apply pid to cgroup
	Apply(path string, pid int) error
	// Remove rm cgroup
	Remove(path string) error
}

var (
	Ins = []Subsystem{
		&CpusetSubSystem{},
		&MemorySubSystem{},
		&CpuSubSystem{},
	}
)
