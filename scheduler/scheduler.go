package scheduler

import (
	"context"
	"github.com/sirupsen/logrus"
	"gpu-jobs-opt/provider"
	"gpu-jobs-opt/types"
	"sync"
)

type Scheduler struct {
	mu       sync.RWMutex
	provider provider.GPUProvider
	//	queue     *JobQueue
	jobs map[string]*types.Job
	//	algorithm Algorithm
	//	predictor PredictorClient
	logger *logrus.Logger
	ctx    context.Context
	cancel context.CancelFunc
}

// Config for scheduler
type Config struct {
	Algorithm    string
	PredictorURL string
}
