package services

import (
	globaltypes "app/agent_grid/internal/global_types"
	"app/agent_grid/internal/response"
	"time"

	"github.com/gorilla/websocket"
)

// ServiceHealthMethods defines the interface for health service
type ServiceHealthMethods interface {
	GetHealth() *HealthCheckResponse
	GetReadiness() *ReadinessResponse
	GetOntologySummary(conn *websocket.Conn, req *globaltypes.RequestOntologySummary) *globaltypes.ResponseOntologySummary
}

type ServiceHealth struct {
	Access *ServiceAccess
}

func NewServiceHealth(access *ServiceAccess) ServiceHealthMethods {
	return &ServiceHealth{
		Access: access,
	}
}

// GetHealth returns the current health status
func (s *ServiceHealth) GetHealth() *HealthCheckResponse {
	return &HealthCheckResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
	}
}

// GetReadiness returns the readiness status
func (s *ServiceHealth) GetReadiness() *ReadinessResponse {
	// TODO: Add database and external service connectivity checks here
	return &ReadinessResponse{
		Ready:     true,
		Timestamp: time.Now(),
	}
}

// GetOntologySummary returns a dummy ontology summary
func (s *ServiceHealth) GetOntologySummary(conn *websocket.Conn, req *globaltypes.RequestOntologySummary) *globaltypes.ResponseOntologySummary {
	response.SendWebSocketMessage(conn, response.WebsocketMessageTypeInfo, "Getting LLM Response...", nil)

	// Simulate processing time
	time.Sleep(4 * time.Second)

	// Return a dummy response for now
	// md format payload
	payload := "## Ontology Summary\n\n- **Entity Count**: 100\n- **Edge Count**: 200\n- **Last Updated**: " + time.Now().Format(time.RFC1123)
	response.SendWebSocketMessage(conn, response.WebsocketMessageTypeLastUpdate, "Last Response", payload)

	return nil
}
