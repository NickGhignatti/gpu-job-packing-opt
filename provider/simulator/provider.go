package simulator

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"gpu-jobs-opt/provider"
	"gpu-jobs-opt/types"
)

// FakeProvider simulates GPU hardware
type FakeProvider struct {
	mu      sync.RWMutex
	devices []*VirtualGPU
	config  Config
	rand    *rand.Rand
}

// Config for simulator
type Config struct {
	NumGPUs       int
	MemoryPerGPU  uint64
	GPUName       string
	RealisticMode bool
}

// VirtualGPU represents a simulated GPU
type VirtualGPU struct {
	mu              sync.RWMutex
	ID              int
	Name            string
	MemoryTotal     uint64
	MemoryUsed      uint64
	BaseUtilization float64
	RunningJobs     map[string]*SimulatedJob
	Metrics         *types.GPUMetrics
}

// SimulatedJob tracks a job on GPU
type SimulatedJob struct {
	JobID            string
	MemoryAllocated  uint64
	ComputeIntensity float64
	StartTime        time.Time
	EstimatedRuntime float64
}

// Register with provider registry
func init() {
	provider.Register("simulator", func() provider.GPUProvider {
		return &FakeProvider{}
	})
}

// Initialize sets up the simulator
func (s *FakeProvider) Initialize(ctx context.Context, config map[string]interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Default configuration
	s.config = Config{
		NumGPUs:       4,
		MemoryPerGPU:  24 * 1024 * 1024 * 1024, // 24GB
		GPUName:       "Simulated GPU A100",
		RealisticMode: true,
	}

	// Override with provided config
	if numGPUs, ok := config["num_gpus"].(int); ok {
		s.config.NumGPUs = numGPUs
	}
	if memGB, ok := config["memory_per_gpu_gb"].(int); ok {
		s.config.MemoryPerGPU = uint64(memGB) * 1024 * 1024 * 1024
	}

	// Initialize random source for noise
	s.rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	// Create virtual GPUs
	s.devices = make([]*VirtualGPU, s.config.NumGPUs)
	for i := 0; i < s.config.NumGPUs; i++ {
		s.devices[i] = &VirtualGPU{
			ID:          i,
			Name:        fmt.Sprintf("%s-%d", s.config.GPUName, i),
			MemoryTotal: s.config.MemoryPerGPU,
			MemoryUsed:  0,
			RunningJobs: make(map[string]*SimulatedJob),
			Metrics:     s.createInitialMetrics(i),
		}
	}

	// Start background metrics updater
	go s.metricsUpdater(ctx)

	return nil
}

// GetDeviceCount returns number of GPUs
func (s *FakeProvider) GetDeviceCount() (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.devices), nil
}

// GetDeviceInfo returns info for one GPU
func (s *FakeProvider) GetDeviceInfo(deviceID int) (*types.GPUInfo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if deviceID < 0 || deviceID >= len(s.devices) {
		return nil, provider.ErrDeviceNotFound
	}

	device := s.devices[deviceID]
	device.mu.RLock()
	defer device.mu.RUnlock()

	// Collect running job IDs
	runningJobIDs := make([]string, 0, len(device.RunningJobs))
	for jobID := range device.RunningJobs {
		runningJobIDs = append(runningJobIDs, jobID)
	}

	memory := types.MemoryMetrics{
		TotalMemoryMB: uint32(device.MemoryTotal),
		UsedMemoryMB:  uint32(device.MemoryUsed),
		FreeMemoryMB:  uint32(device.MemoryTotal - device.MemoryUsed),
	}

	return &types.GPUInfo{
		ID:          uint8(deviceID),
		Name:        device.Name,
		Memory:      memory,
		Utilization: float32(device.BaseUtilization),
		Vendor:      "simulator",
		Capabilities: types.GPUCapabilities{
			SupportMPS:           true,
			SupportMIG:           false,
			SupportMultiProcess:  true,
			MaxConcurrentProcess: 8,
		},
		RunningJobs: runningJobIDs,
	}, nil
}

// GetAllDevices returns info for all GPUs
func (s *FakeProvider) GetAllDevices() ([]*types.GPUInfo, error) {
	count, _ := s.GetDeviceCount()
	devices := make([]*types.GPUInfo, count)

	for i := 0; i < count; i++ {
		info, err := s.GetDeviceInfo(i)
		if err != nil {
			return nil, err
		}
		devices[i] = info
	}

	return devices, nil
}

// GetMetrics returns current metrics for a GPU
func (s *FakeProvider) GetMetrics(deviceID int) (*types.GPUMetrics, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if deviceID < 0 || deviceID >= len(s.devices) {
		return nil, provider.ErrDeviceNotFound
	}

	device := s.devices[deviceID]
	device.mu.RLock()
	defer device.mu.RUnlock()

	// Return copy of metrics
	metrics := *device.Metrics
	return &metrics, nil
}

// GetAllMetrics returns metrics for all GPUs
func (s *FakeProvider) GetAllMetrics() ([]*types.GPUMetrics, error) {
	count, _ := s.GetDeviceCount()
	metrics := make([]*types.GPUMetrics, count)

	for i := 0; i < count; i++ {
		m, err := s.GetMetrics(i)
		if err != nil {
			return nil, err
		}
		metrics[i] = m
	}

	return metrics, nil
}

// EnableMPS - always available in simulator
func (s *FakeProvider) EnableMPS(deviceID int) error {
	return nil
}

// DisableMPS - always available in simulator
func (s *FakeProvider) DisableMPS(deviceID int) error {
	return nil
}

// ConfigureMIG - not supported in simulator
func (s *FakeProvider) ConfigureMIG(deviceID int, config interface{}) error {
	return provider.ErrNotSupported
}

// GetCapabilities returns GPU capabilities
func (s *FakeProvider) GetCapabilities(deviceID int) (*types.GPUCapabilities, error) {
	return &types.GPUCapabilities{
		SupportMPS:           true,
		SupportMIG:           false,
		SupportMultiProcess:  true,
		MaxConcurrentProcess: 8,
	}, nil
}

// Name returns provider name
func (s *FakeProvider) Name() string {
	return "simulator"
}

// Vendor returns vendor name
func (s *FakeProvider) Vendor() string {
	return "simulator"
}

// Close cleans up resources
func (s *FakeProvider) Close() error {
	return nil
}

// AllocateJob simulates allocating a job to a GPU
func (s *FakeProvider) AllocateJob(deviceID int, jobID string, memoryMB int, computeIntensity float64, runtime float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if deviceID < 0 || deviceID >= len(s.devices) {
		return provider.ErrDeviceNotFound
	}

	device := s.devices[deviceID]
	device.mu.Lock()
	defer device.mu.Unlock()

	memoryBytes := uint64(memoryMB) * 1024 * 1024

	// Check if enough memory
	if device.MemoryUsed+memoryBytes > device.MemoryTotal {
		return fmt.Errorf("insufficient GPU memory")
	}

	// Add job
	device.RunningJobs[jobID] = &SimulatedJob{
		JobID:            jobID,
		MemoryAllocated:  memoryBytes,
		ComputeIntensity: computeIntensity,
		StartTime:        time.Now(),
		EstimatedRuntime: runtime,
	}

	device.MemoryUsed += memoryBytes

	return nil
}

// ReleaseJob removes a job from a GPU
func (s *FakeProvider) ReleaseJob(deviceID int, jobID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if deviceID < 0 || deviceID >= len(s.devices) {
		return provider.ErrDeviceNotFound
	}

	device := s.devices[deviceID]
	device.mu.Lock()
	defer device.mu.Unlock()

	job, exists := device.RunningJobs[jobID]
	if !exists {
		return fmt.Errorf("job not found on device")
	}

	device.MemoryUsed -= job.MemoryAllocated
	delete(device.RunningJobs, jobID)

	return nil
}

// createInitialMetrics creates initial metrics
func (s *FakeProvider) createInitialMetrics(deviceID int) *types.GPUMetrics {
	return &types.GPUMetrics{
		GPUID:     uint8(deviceID),
		Timestamp: time.Now(),
		Utilization: types.UtilizationMetrics{
			GPU:    0,
			Memory: 0,
		},
		Memory: types.MemoryMetrics{
			UsedMemoryMB:  0,
			FreeMemoryMB:  uint32(s.config.MemoryPerGPU),
			TotalMemoryMB: uint32(s.config.MemoryPerGPU),
		},
		Temperature: float32(35 + s.rand.Intn(10)),
		PowerUsage:  50.0 + s.rand.Float32()*50.0,
		Processes:   []types.ProcessInfo{},
	}
}

// metricsUpdater updates metrics periodically
func (s *FakeProvider) metricsUpdater(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.updateMetrics()
		}
	}
}

// updateMetrics calculates realistic metrics
func (s *FakeProvider) updateMetrics() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, device := range s.devices {
		device.mu.Lock()

		// Calculate utilization from running jobs
		totalComputeIntensity := 0.0
		processes := make([]types.ProcessInfo, 0, len(device.RunningJobs))

		for jobID, job := range device.RunningJobs {
			totalComputeIntensity += job.ComputeIntensity
			processes = append(processes, types.ProcessInfo{
				PID:        uint8(s.rand.Intn(65535)),
				UsedMemory: float32(job.MemoryAllocated),
				JobID:      jobID,
			})
		}

		// Add realistic noise
		gpuUtil := totalComputeIntensity * 100
		if s.config.RealisticMode {
			gpuUtil += (s.rand.Float64() - 0.5) * 10 // ±5% noise
		}
		gpuUtil = clamp(gpuUtil, 0, 100)

		memoryUtil := float64(device.MemoryUsed) / float64(device.MemoryTotal) * 100

		// Temperature based on utilization
		baseTemp := 35.0
		tempIncrease := gpuUtil * 0.5 // Up to 50°C increase
		temperature := int(baseTemp + tempIncrease)

		// Power based on utilization
		basePower := 50.0
		powerIncrease := gpuUtil * 3.0 // Up to 300W increase
		power := basePower + powerIncrease

		// Update metrics
		device.Metrics = &types.GPUMetrics{
			GPUID:     uint8(device.ID),
			Timestamp: time.Now(),
			Utilization: types.UtilizationMetrics{
				GPU:    float32(gpuUtil),
				Memory: float32(memoryUtil),
			},
			Memory: types.MemoryMetrics{
				UsedMemoryMB:  uint32(device.MemoryUsed),
				FreeMemoryMB:  uint32(device.MemoryTotal - device.MemoryUsed),
				TotalMemoryMB: uint32(device.MemoryTotal),
			},
			Temperature: float32(temperature),
			PowerUsage:  float32(power),
			Processes:   processes,
		}

		device.BaseUtilization = gpuUtil
		device.mu.Unlock()
	}
}

// clamp restricts value to range
func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
