package services

import (
	"time"
)

// ServiceHealthMethods defines the interface for health service
type ServiceHealthMethods interface {
	GetHealth() *HealthCheckResponse
	GetReadiness() *ReadinessResponse
}

// HealthCheckResponse contains health information
type HealthCheckResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
}

// ReadinessResponse contains readiness information
type ReadinessResponse struct {
	Ready     bool      `json:"ready"`
	Timestamp time.Time `json:"timestamp"`
}

// healthService implements ServiceHealthMethods
type healthService struct {
	version string
}

// NewHealthService creates a new health service
func NewHealthService(version string) ServiceHealthMethods {
	return &healthService{
		version: version,
	}
}

// GetHealth returns the current health status
func (s *healthService) GetHealth() *HealthCheckResponse {
	return &HealthCheckResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   s.version,
	}
}

// GetReadiness returns the readiness status
func (s *healthService) GetReadiness() *ReadinessResponse {
	// TODO: Add database and external service connectivity checks here
	return &ReadinessResponse{
		Ready:     true,
		Timestamp: time.Now(),
	}
}
