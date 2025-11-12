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
	UUID       string `json:"uuid"`
	Name       string `json:"name"`
	Image      string
	Commands   []string          `json:"commands"`
	Enviroment map[string]string `json:"enviroment"`
	//	Resource ResourceRequirements
	//  Metadata Metadata
	Status     JobStatus `json:"status"`
	SubmitTime time.Time `json:"submit_time"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	//	Placement *Placement // -> where the job it's scheduled
	Usage *ResourceUsage // -> actualy resource consumption
}

type ResourceUsage struct {
	PeakGPUMemoryMB int32
	AVGGPUUtil      float32
	Runtime         float64
	ExitCode        int8
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
