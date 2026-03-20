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
	Action string `json:"action"`

	// Generic filters (query params)
	Filters map[string]interface{} `json:"filters,omitempty"`

	// Pagination support
	Pagination *utils.Pagination `json:"pagination,omitempty"`

	// Generic payload (for create/update)
	Data any `json:"data,omitempty"`

	// Common identifiers
	Service   string `json:"service,omitempty"`
	FeatureID string `json:"feature_id,omitempty"`
	EntityID  string `json:"entity_id,omitempty"`
	ServiceID string `json:"service_id,omitempty"`
	TeamID    string `json:"team_id,omitempty"`
	APIID     string `json:"api_id,omitempty"`

	// Extra fields
	Role string `json:"role,omitempty"`
	Sort string `json:"sort,omitempty"`
}

type AgentClientMethods interface {
	GenerateResponse(sysPrompt string) (*ResponseAgentModel, error)
	BuildSysPrompt(userPrompt string) (string, error)
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
	Response string `json:"response"` // ⚠️ stringified JSON
}

func (c *OllamaAgentClient) GenerateResponse(sysPrompt string) (*ResponseAgentModel, error) {
	logger := c.Access.logger
	reqBody := OllamaRequest{
		Model:  "tinyllama", // or mistral / phi
		Prompt: sysPrompt,
		Stream: false,
	}

	bodyBytes, _ := json.Marshal(reqBody)

	resp, err := http.Post(c.Access.cfg.Agents.Ollama.Url, "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return nil, fmt.Errorf("internal metrics API returned non-200 status: %d, failed to read body: %w", resp.StatusCode, readErr)
		}
		return nil, fmt.Errorf("internal metrics API returned non-200 status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode internal metrics API response: %w", err)
	}

	bytes, err := json.Marshal(data["response"])
	if err != nil {
		logger.Errorf("failed to marshal internal metrics data: %v", err)
		return nil, fmt.Errorf("failed to marshal internal metrics data: %w", err)
	}

	fmt.Println(string(bytes))
	// Step 1: clean JSON (handles extra text too)
	clean := ExtractJSON(fmt.Sprint(data["response"]))
	if clean == "" {
		return &ResponseAgentModel{
			Data: string(bytes),
		}, nil
	}
	fmt.Println("Cleaned JSON:", clean)
	// Step 2: unmarshal
	var decision ResponseAgentModel
	if err := json.Unmarshal([]byte(clean), &decision); err != nil {
		return nil, fmt.Errorf("failed to parse decision JSON: %w", err)
	}

	fmt.Println(decision)

	return &decision, nil
}

func ExtractJSON(s string) string {
	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")

	if start == -1 || end == -1 || end <= start {
		return ""
	}

	return s[start : end+1]
}

func (c *OllamaAgentClient) BuildSysPrompt(userPrompt string) (string, error) {
	var sb strings.Builder

	sb.WriteString(`
You are an AI agent.

Your task is to determine the user's intent and respond *only* with a JSON object indicating the action to perform and any necessary parameters.
Do NOT include any conversational text or markdown.
Your response MUST be valid JSON.

Rules:
- filters should always be inside "filters" object
- pagination should be inside "pagination" object with page and limit
`)

	sb.WriteString("\n\nAvailable actions:\n")
	i := 0
	for _, action := range agent_config.OntologyAgentActionsList {
		i += 1
		sb.WriteString(fmt.Sprintf("\n%d. **%s**: %s\n", i+1, action.Title, action.Description))

		// Examples
		for _, ex := range action.UserExamples {
			sb.WriteString(fmt.Sprintf("   - \"%s\"\n", ex))
		}

		// Filters
		if len(action.Filters) > 0 {
			sb.WriteString("   Filters:\n")
			for _, f := range action.Filters {
				sb.WriteString(fmt.Sprintf("     - %s (%s): %s\n", f.Key, f.Type, f.Description))
			}
		}

		// Pagination
		if action.Pagination {
			sb.WriteString("   Pagination: supports page and limit\n")
		}

		// JSON format
		sb.WriteString(fmt.Sprintf("   JSON format: %s\n", action.ResponseJSON))

		// API
		sb.WriteString(fmt.Sprintf("   API: [%s] %s\n", action.API.Method, action.API.Endpoint))

		if len(action.API.Headers) > 0 {
			sb.WriteString("   Headers:\n")
			for k, v := range action.API.Headers {
				sb.WriteString(fmt.Sprintf("     - %s: %s\n", k, v))
			}
		}
	}

	sb.WriteString(`
If no action matches:
{"action":"none"}
`)

	sb.WriteString(fmt.Sprintf("\nUser: %s", userPrompt))

	return sb.String(), nil
}
