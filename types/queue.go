package types

import (
	"sync"
)

// JobQueue manages pending jobs (FIFO)
type JobQueue struct {
	mu   sync.RWMutex
	jobs []*Job
}

// NewJobQueue creates a new job queue
func NewJobQueue() *JobQueue {
	return &JobQueue{
		jobs: make([]*Job, 0),
	}
}

// Enqueue adds a job to the queue
func (q *JobQueue) Enqueue(job *Job) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.jobs = append(q.jobs, job)
}

// Dequeue removes and returns the first job
func (q *JobQueue) Dequeue() *Job {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.jobs) == 0 {
		return nil
	}

	job := q.jobs[0]
	q.jobs = q.jobs[1:]
	return job
}

// Peek returns the first job without removing it
func (q *JobQueue) Peek() *Job {
	q.mu.RLock()
	defer q.mu.RUnlock()

	if len(q.jobs) == 0 {
		return nil
	}

	return q.jobs[0]
}

// Len returns the queue length
func (q *JobQueue) Len() int {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return len(q.jobs)
}
