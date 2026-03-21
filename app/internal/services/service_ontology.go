package services

import (
	"app/agent_grid/internal/agent_config"
	"app/agent_grid/internal/clients"
	"app/agent_grid/internal/config"
	globaltypes "app/agent_grid/internal/global_types"
	"app/agent_grid/internal/response"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
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

	agentConfig, modelConfig, ok := config.GetAgentForFeature(s.Access.Cfg, "OntologyAgent")
	if !ok {
		logger.Errorf("Failed to get agent for feature: OntologyAgent")
		response.SendWebSocketMessage(conn, response.WebsocketMessageTypeError, "Agent configuration not found", nil)
		return nil
	}

	logger.Infof("Using model: %s", agentConfig.Model)
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
	plannerResp, err := s.Access.Clients.ClientAgent.GenerateResponse(clients.AgentName(modelConfig.Agent), agentConfig, plannerPrompt)
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
	var stepResults = map[int]any{}

	for _, step := range plan.Plan {
		var count = 0
		action, ok := agent_config.OntologyAgentActionsList[agent_config.ActionName(step.Action)]
		if !ok {
			continue
		}

		type callStack struct {
			endpoint string
		}

		var stack []callStack

		// =========================
		// 🔥 RESOLVE PARAMS (CHAINING)
		// =========================
		params := s.ResolveParams(mergeMaps(step.PathParams, step.QueryParams), stepResults)
		step.QueryParams = params
		step.Body = s.ResolveParams(step.Body, stepResults)
		step.PathParams = params
		queryParams := buildQueryParams(params)

		// =========================
		// 🔥 BUILD ENDPOINT
		// =========================
		endpoint := action.API.Endpoint
		for k, v := range step.PathParams {
			pathParamName := "{" + k + "}"
			if !strings.Contains(endpoint, pathParamName) {
				continue
			}

			// if v is a an array
			if arr, ok := v.([]interface{}); ok {
				for _, item := range arr {
					stack = append(stack, callStack{
						endpoint: strings.ReplaceAll(endpoint, pathParamName, fmt.Sprint(item)),
					})
				}

			} else {
				stack = append(stack, callStack{
					endpoint: strings.ReplaceAll(endpoint, pathParamName, fmt.Sprint(v)),
				})
			}
			break // only replace first occurrence as multiple path params are not supported in this version
		}

		if len(stack) == 0 {
			stack = append(stack, callStack{
				endpoint: endpoint,
			})
		}

		var resultApiResp any
		for _, item := range stack {
			endpoint = item.endpoint

			action.API.Endpoint = endpoint + queryParams
			step.Endpoint = endpoint + queryParams
			response.SendWebSocketMessage(conn,
				response.WebsocketMessageTypeInfo,
				fmt.Sprintf("Executing: %s", step.Action),
				step,
			)

			// =========================
			// 🔥 EXECUTE API
			// =========================
			count += 1
			apiResp, err := s.Access.Clients.ClientGeneral.SendRequestByAgent(
				context.Background(),
				&action,
				action.API.Headers,
				step.Body,
				step.QueryParams,
			)

			if err != nil {
				logger.Error("API error: %v", err)
				continue
			}

			resultApiResp, err = mergeResults(resultApiResp, apiResp)
			if err != nil {
				logger.Error("Merge results error: %v", err)
				resultApiResp = apiResp // fallback to latest result if merge fails
			}
		}

		// store result for chaining
		stepResults[step.Step] = resultApiResp

		response.SendWebSocketMessage(conn,
			response.WebsocketMessageTypeInfo,
			fmt.Sprintf("Step %d done with %d API calls", step.Step, count),
			resultApiResp,
		)
	}

	// =========================
	// 4. FINAL SUMMARY
	// =========================
	var stepResultsStr []string
	for _, res := range stepResults {
		stepResultsStr = append(stepResultsStr, fmt.Sprint(res))
	}

	finalPrompt := agent_config.BuildSummaryPrompt(req.Prompt, stepResultsStr)
	response.SendWebSocketMessage(conn,
		response.WebsocketMessageTypeInfo,
		"Summarising final response...",
		nil,
	)

	finalResp, err := s.Access.Clients.ClientAgent.GenerateResponse(clients.AgentName(modelConfig.Agent), agentConfig, finalPrompt)
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

func (s *ServiceHealth) ResolveParams(params map[string]any, stepResults map[int]any) map[string]any {
	logger := s.Access.Logger

	if params == nil {
		return nil
	}

	resolved := map[string]any{}

	for k, v := range params {
		if v == nil || fmt.Sprint(v) == "" {
			continue
		}
		strVal := fmt.Sprint(v)

		// =========================
		// HANDLE STEP REFERENCES
		// =========================
		if strings.HasPrefix(strVal, "$") {

			// normalize $step_1 → $1
			strVal = strings.ReplaceAll(strVal, "$step_", "$")

			// split: $1.result.id[0]
			parts := strings.Split(strVal, ".")
			if len(parts) < 3 {
				logger.Errorf("ResolveParams: invalid reference format | key=%s value=%s", k, strVal)
				resolved[k] = v
				continue
			}

			// -------------------------
			// STEP NUMBER
			// -------------------------
			stepStr := strings.TrimPrefix(parts[0], "$")
			stepNum, err := strconv.Atoi(stepStr)
			if err != nil {
				logger.Errorf("ResolveParams: invalid step number | key=%s value=%s err=%v", k, strVal, err)
				resolved[k] = v
				continue
			}

			prev, ok := stepResults[stepNum]
			if !ok {
				logger.Errorf("ResolveParams: step result not found | step=%d key=%s", stepNum, k)
				resolved[k] = v
				continue
			}

			// -------------------------
			// PARSE prev → map
			// -------------------------
			var prevMap map[string]interface{}

			switch val := prev.(type) {
			case string:
				if err := json.Unmarshal([]byte(val), &prevMap); err != nil {
					logger.Errorf("ResolveParams: failed to unmarshal string | step=%d err=%v raw=%s", stepNum, err, val)
					resolved[k] = v
					continue
				}
			default:
				b, err := json.Marshal(val)
				if err != nil {
					logger.Errorf("ResolveParams: marshal failed | step=%d err=%v", stepNum, err)
					resolved[k] = v
					continue
				}
				if err := json.Unmarshal(b, &prevMap); err != nil {
					logger.Errorf("ResolveParams: unmarshal failed | step=%d err=%v", stepNum, err)
					resolved[k] = v
					continue
				}
			}

			// -------------------------
			// GET RESULT
			// -------------------------
			resultRaw, ok := prevMap["result"]
			if !ok {
				resultRaw, ok = prevMap["data"]
				if !ok {
					logger.Errorf("ResolveParams: 'result' key missing | step=%d data=%v", stepNum, prevMap)
					resolved[k] = v
					continue
				}
			}

			var resultArr []interface{}

			switch r := resultRaw.(type) {
			case []interface{}:
				resultArr = r
			default:
				resultArr = []interface{}{r} // wrap single object
			}

			field := parts[2]

			// =========================
			// HANDLE id[0] ONLY
			// =========================
			if strings.Contains(field, "[") && !strings.HasSuffix(field, "[]") {

				fieldName := field[:strings.Index(field, "[")]
				indexStr := field[strings.Index(field, "[")+1 : strings.Index(field, "]")]

				idx, err := strconv.Atoi(indexStr)
				if err != nil {
					logger.Errorf("ResolveParams: invalid index | field=%s err=%v", field, err)
					resolved[k] = v
					continue
				}

				if len(resultArr) <= idx {
					logger.Errorf("ResolveParams: index out of bounds | idx=%d len=%d", idx, len(resultArr))
					resolved[k] = v
					continue
				}

				itemBytes, err := json.Marshal(resultArr[idx])
				if err != nil {
					logger.Errorf("ResolveParams: marshal item failed | err=%v", err)
					resolved[k] = v
					continue
				}

				var obj map[string]interface{}
				if err := json.Unmarshal(itemBytes, &obj); err != nil {
					logger.Errorf("ResolveParams: unmarshal item failed | err=%v", err)
					resolved[k] = v
					continue
				}

				resolved[k] = obj[fieldName]
				continue
			}

			// =========================
			// HANDLE id[] (ALL VALUES)
			// =========================
			if strings.HasSuffix(field, "[]") {

				fieldName := strings.TrimSuffix(field, "[]")
				var values []interface{}

				for i, item := range resultArr {
					itemBytes, err := json.Marshal(item)
					if err != nil {
						logger.Errorf("ResolveParams: marshal item failed | idx=%d err=%v", i, err)
						continue
					}

					var obj map[string]interface{}
					if err := json.Unmarshal(itemBytes, &obj); err != nil {
						logger.Errorf("ResolveParams: unmarshal item failed | idx=%d err=%v", i, err)
						continue
					}

					values = append(values, obj[fieldName])
				}

				resolved[k] = values
				continue
			}

			// =========================
			// DEFAULT: id → first item
			// =========================
			if len(resultArr) == 0 {
				logger.Errorf("ResolveParams: empty result array | step=%d", stepNum)
				resolved[k] = v
				continue
			}

			itemBytes, err := json.Marshal(resultArr[0])
			if err != nil {
				logger.Errorf("ResolveParams: marshal first item failed | err=%v", err)
				resolved[k] = v
				continue
			}

			var obj map[string]interface{}
			if err := json.Unmarshal(itemBytes, &obj); err != nil {
				logger.Errorf("ResolveParams: unmarshal first item failed | err=%v", err)
				resolved[k] = v
				continue
			}

			resolved[k] = obj[field]

		} else {
			// -------------------------
			// NORMAL VALUE
			// -------------------------
			resolved[k] = v
		}
	}

	return resolved
}

func mergeMaps(m1, m2 map[string]any) map[string]any {
	out := map[string]any{}
	for k, v := range m1 {
		out[k] = v
	}
	for k, v := range m2 {
		out[k] = v
	}
	return out
}

func copyMap(m map[string]any) map[string]any {
	if m == nil {
		return nil
	}
	out := make(map[string]any)
	for k, v := range m {
		out[k] = v
	}
	return out
}

func mergeResults(r1, r2 any) (any, error) {
	m1, err := toMap(r1)
	if err != nil {
		return nil, err
	}

	if m1 == nil {
		m1 = map[string]any{}
	}

	m2, err := toMap(r2)
	if err != nil {
		return nil, err
	}

	if m2 == nil {
		m2 = map[string]any{}
	}

	var merged []any

	if r, ok := m1["result"].([]any); ok {
		merged = append(merged, r...)
	}

	if r, ok := m2["result"].([]any); ok {
		merged = append(merged, r...)
	}

	m1["result"] = merged
	return m1, nil
}

func buildQueryParams(data map[string]any) string {
	params := url.Values{}
	a := false

	for key, val := range data {
		switch v := val.(type) {

		// handle array → multiple params
		case []any:
			for _, item := range v {
				params.Add(key, fmt.Sprintf("%v", item))
				a = true
			}

		// handle normal value
		default:
			params.Add(key, fmt.Sprintf("%v", v))
			a = true
		}
	}

	if !a {
		return ""
	}

	return "?" + params.Encode()
}

func toMap(input any) (map[string]any, error) {
	var m = make(map[string]any)

	b, err := json.Marshal(input)
	if err != nil {
		return m, err
	}

	// first attempt: normal unmarshal
	if err := json.Unmarshal(b, &m); err == nil {
		return m, nil
	}

	// fallback: input is JSON string → unquote → unmarshal again
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return m, err
	}

	if err := json.Unmarshal([]byte(str), &m); err != nil {
		return m, err
	}

	return m, nil
}
