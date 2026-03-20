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
	time.Sleep(2 * time.Second)

	response.SendWebSocketMessage(conn, response.WebsocketMessageTypeUpdate, "Hitting API response...", nil)

	// Simulate processing time
	time.Sleep(2 * time.Second)

	response.SendWebSocketMessage(conn, response.WebsocketMessageTypeUpdate, "Getting LLM Response...", nil)

	// Simulate processing time
	time.Sleep(4 * time.Second)

	// Return a dummy response for now
	// md format payload
	payload := "## Ontology Summary\n\n- **Entity Count**: 100\n- **Edge Count**: 200\n- **Last Updated**: " + time.Now().Format(time.RFC1123)
	payload = "**Feature List Response**\n======================\n\n### Successful API Request\n\nYour request to retrieve the feature list was successful!\n\n### Key Features\n\n* **kyc_status**: This is one of the features available in your system. Here are some key details about this feature:\n\n    * **Code:** kyc_status\n    * **Name:** kyc_status\n    * **Description:** auto-created by ontology consumer\n    * **Status:** Active (is_active: true)\n\n### Summary\n\nIn summary, you have one feature available in your system, which is the `kyc_status` feature. This feature is active and was created automatically by the ontology consumer.\n\n### API Response Details\n\nHere are some additional details about the API response:\n\n* **Code:** 00000 (successful request)\n* **Message:** success\n* **Pagination:**\n    + Current page: 1\n    + Page size: 10\n    + Total count: 1\n    + Total pages: 1"
	response.SendWebSocketMessage(conn, response.WebsocketMessageTypeLastUpdate, "Last Response", payload)

	return nil
}
