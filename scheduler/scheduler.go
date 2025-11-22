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
	queue    *types.JobQueue
	jobs     map[string]*types.Job
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

func New(provider provider.GPUProvider, config Config) (*Scheduler, error) {
	ctx, cancel := context.WithCancel(context.Background())

	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	// TODO : algorithm selection (kinda best-fit or first-fit)
	
	s := &Scheduler{
		mu:       sync.RWMutex{},
		provider: provider,
		queue:    types.NewJobQueue(),
		jobs:     make(map[string]*types.Job),
		logger:   logger,
		ctx:      ctx,
		cancel:   cancel,
	}

	if config.PredictorURL != "" {
		// create predictor client if provided
	}

	return s, nil
}

func (s *Scheduler) SubmitJob(job *types.Job) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// TODO : should predict here

	s.jobs[job.UUID] = job
	s.queue.Enqueue(job)
	s.logger.Infof("Job %s submitted (%s)", job.UUID, job.Name)

	return nil
}
