package types

import "time"

type GPUInfo struct {
	ID           uint8           `json:"id"`
	Name         string          `json:"name"`
	Memory       MemoryMetrics   `json:"memory"`
	Utilization  float32         `json:"utilization"`
	Temperatire  float32         `json:"temperatire"`
	PowerUsage   float32         `json:"power_usage"`
	Vendor       string          `json:"vendor"`
	Capabilities GPUCapabilities `json:"capabilities"`
	RunningJobs  []string        `json:"running_jobs"`
}

type GPUCapabilities struct {
	SupportMPS           bool   `json:"support_mps"`
	SupportMIG           bool   `json:"support_mig"`
	SupportMultiProcess  bool   `json:"support_multi_process"`
	MaxConcurrentProcess uint64 `json:"max_concurrent_process"`
}

type GPUMetrics struct {
	GPUID       uint8              `json:"gpuid"`
	Timestamp   time.Time          `json:"timestamp"`
	Utilization UtilizationMetrics `json:"utilization"`
	Memory      MemoryMetrics      `json:"memory"`
	Temperatire float32            `json:"temperatire"`
	PowerUsage  float32            `json:"power_usage"`
	Processes   []ProcessInfo      `json:"processes"`
}

type UtilizationMetrics struct {
	GPU    float32 `json:"gpu"`
	Memory float32 `json:"memory"`
}

type MemoryMetrics struct {
	TotalMemoryMB uint32 `json:"total_memory_mb"`
	UsedMemoryMB  uint32 `json:"used_memory_mb"`
	FreeMemoryMB  uint32 `json:"free_memory_mb"`
}

type ProcessInfo struct {
	PID        uint8   `json:"pid"`
	UsedMemory float32 `json:"used_memory"`
	Name       string  `json:"name"`
	JobID      string  `json:"job_id"`
}
