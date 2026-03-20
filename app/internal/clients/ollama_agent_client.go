package clients

import (
	"app/agent_grid/internal/agent_config"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"bitbucket.org/fyscal/be-commons/pkg/utils"
)

type ResponseAgentModel struct {
	Action string `json:"action,omitempty"`

	// Planner support
	Plan []agent_config.PlanStep `json:"plan,omitempty"`

	// Generic filters (query params)
	Filters map[string]interface{} `json:"filters,omitempty"`

	// Pagination
	Pagination *utils.Pagination `json:"pagination,omitempty"`

	// Generic payload
	Data any `json:"data,omitempty"`

	// IDs
	ServiceID string `json:"service_id,omitempty"`
	FeatureID string `json:"feature_id,omitempty"`
	EntityID  string `json:"entity_id,omitempty"`
	TeamID    string `json:"team_id,omitempty"`
	APIID     string `json:"api_id,omitempty"`

	// Extra
	Role string `json:"role,omitempty"`
	Sort string `json:"sort,omitempty"`
}

type AgentClientMethods interface {
	GenerateResponse(prompt string) (*ResponseAgentModel, error)
	BuildPlannerPrompt(userPrompt string, actions []agent_config.LLMAction) string
}

type OllamaAgentClient struct {
	Access *clientAccess
}

func NewOllamaAgentClient(access *clientAccess) AgentClientMethods {
	return &OllamaAgentClient{
		Access: access,
	}
}

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaResponse struct {
	Response string `json:"response"`
}

// =========================
// 🔥 GENERATE RESPONSE (GENERIC)
// =========================
func (c *OllamaAgentClient) GenerateResponse(prompt string) (*ResponseAgentModel, error) {
	logger := c.Access.logger

	reqBody := OllamaRequest{
		Model:  "llama3",
		Prompt: prompt,
		Stream: false,
	}

	bodyBytes, _ := json.Marshal(reqBody)

	resp, err := http.Post(c.Access.cfg.Agents.Ollama.Url, "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	rawBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ollama error: %d %s", resp.StatusCode, string(rawBody))
	}

	var data map[string]interface{}
	if err := json.Unmarshal(rawBody, &data); err != nil {
		return nil, err
	}

	raw := fmt.Sprint(data["response"])

	// 🔥 Extract JSON safely
	clean := ExtractJSON(raw)
	if clean == "" {
		return &ResponseAgentModel{
			Data: raw,
		}, nil
	}

	var result ResponseAgentModel
	if err := json.Unmarshal([]byte(clean), &result); err != nil {
		logger.Errorf("failed to parse LLM JSON: %v", err)

		// fallback → treat as plain text
		return &ResponseAgentModel{
			Data: raw,
		}, nil
	}

	return &result, nil
}

// =========================
// 🔥 PLANNER PROMPT (NEW)
// =========================
func (c *OllamaAgentClient) BuildPlannerPrompt(userPrompt string, actions []agent_config.LLMAction) string {
	actionsJSON, _ := json.Marshal(actions)

	return fmt.Sprintf(`
You are an AI agent that plans API calls.

User Query:
%s

Available Actions:
%s

Rules:
- Only use provided actions
- Prefer minimum steps
- Use multiple steps if debugging/analysis
- Do NOT hallucinate
- Return ONLY JSON

Output format:
{
  "plan": [
    {
      "step": 1,
      "action": "action_name",
      "reason": "why",
      "depends_on": null
    }
  ]
}
`, userPrompt, string(actionsJSON))
}

// =========================
// 🔥 JSON EXTRACTOR
// =========================
func ExtractJSON(s string) string {
	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")

	if start == -1 || end == -1 || end <= start {
		return ""
	}
	return s[start : end+1]
}
