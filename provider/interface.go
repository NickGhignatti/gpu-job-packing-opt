package provider

import (
	"context"
	"gpu-jobs-opt/types"
)

type GPUProvider interface {
	Initialize(ctx context.Context, config map[string]interface{}) error
	GetDeviceCount() (int, error)
	GetDeviceInfo(deviceID int) (*types.GPUInfo, error)
	GetAllDevices() ([]*types.GPUInfo, error)
	GetMetrics(deviceID int) (*types.GPUMetrics, error)
	GetAllMetrics() ([]*types.GPUMetrics, error)
	GetCapabilities(deviceID int) (*types.GPUCapabilities, error)

	Name() string
	Vendor() string
	Close() error
}

// Registry holds all available providers
type Registry struct {
	providers map[string]Factory
}

type Factory func() GPUProvider

type Error struct {
	Message string
	Err     error
}

var (
	ErrProviderNotFound = &Error{Message: "provider not found"}
)

var globalRegistry = &Registry{
	providers: make(map[string]Factory),
}

func Register(name string, factory Factory) {
	globalRegistry.providers[name] = factory
}

func Get(name string) (Factory, error) {
	factory, exists := globalRegistry.providers[name]
	if !exists {
		return nil, ErrProviderNotFound.Err
	}
	return factory, nil
}

func ListProvider() []string {
	keys := make([]string, len(globalRegistry.providers))

	i := 0
	for k := range globalRegistry.providers {
		keys[i] = k
		i++
	}

	return keys
}

func (e *Error) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Err
}
