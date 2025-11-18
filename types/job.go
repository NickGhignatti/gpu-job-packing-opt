package types

import "time"
import "github.com/google/uuid"

type JobStatus int

const (
	Pending JobStatus = iota
	Scheduled
	Running
	Completed
	Failed
	Cancelled
)

type Job struct {
	UUID       string               `json:"uuid"`
	Name       string               `json:"name"`
	Image      string               `json:"image"`
	Commands   []string             `json:"commands"`
	Enviroment map[string]string    `json:"enviroment"`
	Resource   ResourceRequirements `json:"resource"`
	Metadata   JobMetadata          `json:"metadata"`
	Status     JobStatus            `json:"status"`
	SubmitTime time.Time            `json:"submit_time"`
	StartTime  time.Time            `json:"start_time"`
	EndTime    time.Time            `json:"end_time"`
	Placement  *Placement           `json:"placement"`
	Usage      *ResourceUsage       `json:"usage"`
}

type ResourceUsage struct {
	PeakGPUMemoryMB uint32  `json:"peak_gpu_memory_mb"`
	AVGGPUUtil      float32 `json:"avggpu_util"`
	Runtime         float64 `json:"runtime"`
	ExitCode        int8    `json:"exit_code"`
}

type ResourceRequirements struct {
	GPUMemoryMB    uint32  `json:"gpu_memory_mb"`
	GPUCount       uint8   `json:"gpu_count"`
	GPUFraction    float64 `json:"gpu_fraction"`
	CPUCores       uint8   `json:"cpu_cores"`
	SystemMemoryMB uint32  `json:"system_memory_mb"`
}

type JobMetadata struct {
	ModelType      string `json:"model_type"`
	BatchSize      uint16 `json:"batch_size"`
	SequenceLength uint16 `json:"sequence_length"`
	NumParameters  uint32 `json:"num_parameters"`
	// TODO : rn we'll keep complexity level easy
	// MixedPrecision     bool
	// GradientCheckpoint bool
	// HistoricalMemoryMB []uint
	// HistoricalRuntime  []float64
}

type Placement struct {
	GPUID       uint8  `json:"gpuid"`
	NodeID      string `json:"node_id"`
	UseMPS      bool   `json:"use_mps"`
	ContainerID string `json:"container_id"`
}

func NewJob(name, image string, commands []string) *Job {
	return &Job{
		UUID:       uuid.New().String(),
		Name:       name,
		Image:      image,
		Commands:   commands,
		Enviroment: make(map[string]string),
		Status:     Pending,
		SubmitTime: time.Now(),
	}
}
