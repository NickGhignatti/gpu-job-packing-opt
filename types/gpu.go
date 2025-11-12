package types

type GPUInfo struct {
	ID   uint8
	Name string
	// FIXME : better decision making about types
	TotalMemory  uint64
	UsedMemory   uint64
	FreeMemory   uint64
	Utilization  float32
	Temperatire  float32
	PowerUsage   float32
	Vendor       string
	Capabilities GPUCapabilities
	RunningJobs  []string
}

type GPUCapabilities struct {
	SupportMPS           bool
	SupportMIG           bool
	SupportMultiProcess  bool
	MaxConcurrentProcess uint64
}
