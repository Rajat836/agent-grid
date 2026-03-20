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

func (s *ServiceHealth) GetOntologySummary(conn *websocket.Conn, req *globaltypes.RequestOntologySummary) *globaltypes.ResponseOntologySummary {
	logger := s.Access.Logger
	clients := s.Access.Clients

	response.SendWebSocketMessage(conn, response.WebsocketMessageTypeStart, "Thinking...", nil)

	// =========================
	// 1. FILTER ACTIONS
	// =========================
	filteredActions := agent_config.FilterActions(req.Prompt)
	llmActions := agent_config.ToLLMActions(filteredActions)

	// =========================
	// 2. PLANNER PROMPT
	// =========================
	plannerPrompt := agent_config.BuildPlannerPrompt(req.Prompt, llmActions)

	plannerResp, err := clients.ClientOllama.GenerateResponse(plannerPrompt)
	if err != nil {
		logger.Error("Planner error: %v", err)
		response.SendWebSocketMessage(conn, response.WebsocketMessageTypeError, "Planner failed", err.Error())
		return nil
	}

	var plan agent_config.Plan
	a, _ := json.Marshal(plannerResp)
	fmt.Println(string(a))
	err = json.Unmarshal(a, &plan)
	if err != nil {
		logger.Error("Planner error: %v", err)
		response.SendWebSocketMessage(conn, response.WebsocketMessageTypeError, "Failed to parse plan", err.Error())
		return nil
	}

	if len(plan.Plan) == 0 {
		response.SendWebSocketMessage(conn, response.WebsocketMessageTypeError, "No plan generated", nil)
		return nil
	}

	response.SendWebSocketMessage(conn, response.WebsocketMessageTypeInfo, "Plan Generated", plan)

	// =========================
	// 3. EXECUTION LOOP
	// =========================
	var stepResults []any

	for _, step := range plan.Plan {

		action, ok := agent_config.OntologyAgentActionsList[agent_config.ActionName(step.Action)]
		if !ok {
			continue
		}

		response.SendWebSocketMessage(conn, response.WebsocketMessageTypeInfo,
			fmt.Sprintf("Executing: %s", step.Action), step.Reason)

		// NOTE: currently empty params (LLM can be extended later to send params)
		apiResp, err := clients.ClientGeneral.SendRequestByAgent(
			context.Background(),
			&action,
			action.API.Headers,
			nil,
			nil,
		)

		if err != nil {
			logger.Error("API error: %v", err)
			continue
		}

		stepResults = append(stepResults, apiResp)
		response.SendWebSocketMessage(conn,
			response.WebsocketMessageTypeInfo,
			fmt.Sprintf("Step %d done", step.Step),
			apiResp,
		)
	}

	// =========================
	// 4. FINAL SUMMARY
	// =========================
	var stepResultsStr []string
	for _, res := range stepResults {
		stepResultsStr = append(stepResultsStr, fmt.Sprint(res))
	}

	finalPrompt := fmt.Sprintf(`
User Query:
%s

Step Results:
%s

Provide final concise markdown summary.
`, req.Prompt, strings.Join(stepResultsStr, "\n"))

	response.SendWebSocketMessage(conn,
		response.WebsocketMessageTypeInfo,
		"Summarising final response...",
		nil,
	)

	finalResp, err := clients.ClientOllama.GenerateResponse(finalPrompt)
	if err != nil {
		logger.Error("Final summary error: %v", err)
		return nil
	}

	response.SendWebSocketMessage(conn,
		response.WebsocketMessageTypeLastUpdate,
		"Final Response",
		ConvertEscapedToMarkdown(finalResp.Data),
	)

	return nil
}

func ConvertEscapedToMarkdown(input any) string {
	s, err := strconv.Unquote(fmt.Sprint(input))
	if err != nil {
		return fmt.Sprint(input)
	}
	return strings.ReplaceAll(s, "\r\n", "\n")
}
