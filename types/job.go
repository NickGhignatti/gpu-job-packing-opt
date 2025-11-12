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
	PeakGPUMemoryMB uint32
	AVGGPUUtil      float32
	Runtime         float64
	ExitCode        int8
}

type ResourceRequirements struct {
	GPUMemoryMB    uint32
	GPUCount       uint8
	GPUFraction    float64
	CPUCores       uint8
	SystemMemoryMB uint32
}

type JobMetadata struct {
	ModelType      string
	BatchSize      uint16
	SequenceLength uint16
	NumParameters  uint32
	// TODO : rn we'll keep complexity level easy
	// MixedPrecision     bool
	// GradientCheckpoint bool
	// HistoricalMemoryMB []uint
	// HistoricalRuntime  []float64
}

type Placement struct {
	GPUID       uint8
	NodeID      string
	UseMPS      bool
	ContainerID string
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
