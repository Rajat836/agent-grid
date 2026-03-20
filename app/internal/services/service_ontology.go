package services

import (
	"app/agent_grid/internal/agent_config"
	globaltypes "app/agent_grid/internal/global_types"
	"app/agent_grid/internal/response"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
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
	var (
		logger  = s.Access.Logger
		clients = s.Access.Clients
	)

	response.SendWebSocketMessage(conn, response.WebsocketMessageTypeStart, "Getting LLM Response...", nil)

	prompt, _ := clients.ClientOllama.BuildSysPrompt(req.Prompt)
	resp, err := clients.ClientOllama.GenerateResponse(prompt)
	if err != nil {
		logger.Error("Error getting ontology summary from Ollama: %v", err)
		response.SendWebSocketMessage(conn, response.WebsocketMessageTypeError, "Error getting ontology summary", err.Error())
		return nil
	}

	// Send the LLM response back to the client
	response.SendWebSocketMessage(conn, response.WebsocketMessageTypeInfo, "LLM Response", resp)

	action, ok := agent_config.OntologyAgentActionsList[resp.Action]
	if !ok {
		logger.Error("Unknown action from LLM response: %s", resp.Action)
		response.SendWebSocketMessage(conn, response.WebsocketMessageTypeError, "Unknown action in LLM response", resp.Action)
		return nil
	}

	clientResp, err := clients.ClientGeneral.SendRequestByAgent(context.Background(), &action, action.API.Headers, resp.Data, resp.Filters)
	if err != nil {
		logger.Error("Error getting ontology summary from OntologyBot: %v", err)
		response.SendWebSocketMessage(conn, response.WebsocketMessageTypeError, "Error processing ontology summary", err.Error())
		return nil
	}

	newPrompt := fmt.Sprintf(`
You are an AI assistant.
Given the following API response data for %s:
%s
Summarize this information in a concise and user-friendly manner in markdown format.
Highlight key details and any important statuses or outcomes.

User's Instructions for this data:
%s
`, action.Name, clientResp, req.Prompt)

	clientResp2, err := clients.ClientOllama.GenerateResponse(newPrompt)
	if err != nil {
		logger.Error("Error summarizing ontology summary from Ollama: %v", err)
		response.SendWebSocketMessage(conn, response.WebsocketMessageTypeError, "Error summarizing ontology summary", err.Error())
		return nil
	}

	var result string
	err = json.Unmarshal([]byte(fmt.Sprint(clientResp2.Data)), &result)
	if err != nil {
	}

	response.SendWebSocketMessage(conn, response.WebsocketMessageTypeLastUpdate, "OntologyBot Response", ConvertEscapedToMarkdown(clientResp2.Data))
	return nil
}

func ConvertEscapedToMarkdown(input any) string {
	// Step 1: Unquote the string (handles \" and \\n etc.)
	unquoted, err := strconv.Unquote(fmt.Sprint(input))
	if err != nil {
		return err.Error()
	}

	// Step 2: Normalize newlines (optional safety)
	unquoted = strings.ReplaceAll(unquoted, "\r\n", "\n")

	return unquoted
}
