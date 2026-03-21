package agent_config

import (
	"encoding/json"
	"fmt"
	"strings"
)

// =========================
// LLM STRUCTS
// =========================

type LLMAction struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Params       []string `json:"params"`
	UserExamples []string `json:"user_examples,omitempty"`
}

type PlanStep struct {
	Step      int    `json:"step"`
	Action    string `json:"action"`
	Reason    string `json:"reason"`
	DependsOn *int   `json:"depends_on"`

	QueryParams map[string]any `json:"query_params"`
	Body        map[string]any `json:"body"`
	PathParams  map[string]any `json:"path_params"`

	// Just for logging
	Endpoint string `json:"endpoint"`
}
type Plan struct {
	Plan []PlanStep `json:"plan"`
}

// =========================
// DOMAIN MAPPING
// =========================

// var DomainActionsMap = map[string][]ActionName{
// 	"features": {GetFeatures},
// 	"entities": {GetEntities, GetEntityMetrics},
// 	"services": {GetServices, GetServiceDeployments},
// 	"apis":     {GetAPIs, GetAPIMetrics},
// }

// =========================
// FILTERING
// =========================

func DetectDomains(query string) []string {
	q := strings.ToLower(query)

	var domains []string

	if strings.Contains(q, "feature") {
		domains = append(domains, "features")
	}
	if strings.Contains(q, "entity") || strings.Contains(q, "step") {
		domains = append(domains, "entities")
	}
	if strings.Contains(q, "service") || strings.Contains(q, "deploy") {
		domains = append(domains, "services")
	}
	if strings.Contains(q, "api") {
		domains = append(domains, "apis")
	}
	if strings.Contains(q, "team") {
		domains = append(domains, "team")
	}
	if strings.Contains(q, "kpi") {
		domains = append(domains, "kpi")
	}

	if len(domains) == 0 {
		return []string{"features"}
	}

	return domains
}

func createDomainActionMap() map[string][]ActionName {
	domainActionMap := make(map[string][]ActionName)

	for actionName, action := range OntologyAgentActionsList {
		if strings.Contains(strings.ToLower(action.Description), "feature") {
			domainActionMap["features"] = append(domainActionMap["features"], actionName)
		}
		if strings.Contains(strings.ToLower(action.Description), "entity") {
			domainActionMap["entities"] = append(domainActionMap["entities"], actionName)
		}
		if strings.Contains(strings.ToLower(action.Description), "service") {
			domainActionMap["services"] = append(domainActionMap["services"], actionName)
		}
		if strings.Contains(strings.ToLower(action.Description), "api") {
			domainActionMap["apis"] = append(domainActionMap["apis"], actionName)
		}
		if strings.Contains(strings.ToLower(action.Description), "team") {
			domainActionMap["team"] = append(domainActionMap["team"], actionName)
		}
		if strings.Contains(strings.ToLower(action.Description), "kpi") {
			domainActionMap["kpi"] = append(domainActionMap["kpi"], actionName)
		}
	}

	return domainActionMap
}

func FilterActions(query string) []Action {
	domains := DetectDomains(query)
	domainActionMap := createDomainActionMap()

	actionMap := map[ActionName]Action{}

	for _, d := range domains {
		for _, a := range domainActionMap[d] {
			actionMap[a] = OntologyAgentActionsList[a]
		}
	}

	var res []Action
	for _, v := range actionMap {
		res = append(res, v)
	}
	return res
}

// =========================
// LLM ACTION CONVERSION
// =========================

func ToLLMActions(actions []Action) []LLMAction {
	var result []LLMAction

	for _, a := range actions {
		var params []string

		for _, p := range a.PathParams {
			params = append(params, p.Key)
		}
		for _, p := range a.QueryParams {
			params = append(params, p.Key)
		}

		result = append(result, LLMAction{
			Name:         string(a.Name),
			Description:  a.Description,
			UserExamples: a.UserExamples,
			Params:       params,
		})
	}

	return result
}

// =========================
// PLANNER PROMPT
// =========================

func BuildPlannerPrompt(query string, actions []LLMAction) string {
	actionsJSON, _ := json.Marshal(actions)

	return fmt.Sprintf(`
You are an AI agent that plans API calls.

User Query:
%s

Available Actions:
%s

=========================
STRICT RULES (IMPORTANT)
=========================

1. Use ONLY provided actions
2. Return step-by-step plan
3. ALWAYS include params in each step

-------------------------
PARAM USAGE RULES
-------------------------

1. query_params:
   - ONLY for filters
   - Example: ?code=abc

2. path_params:
   - MUST be used when endpoint contains path variables
   - Example:
     Endpoint: /entities/{entity_id}/apis
     THEN:
     "path_params": { "entity_id": 123 }

   - ❌ NEVER put path params inside query_params

3. body:
   - ONLY for POST/PATCH

-------------------------
CHAINING RULES
-------------------------

- Use previous step outputs like:
  "$1.result.id"
  "$2.result.id[]"

- If multiple IDs:
  ALWAYS use [] (array)

-------------------------
EXAMPLES (VERY IMPORTANT)
-------------------------

Example 1:
Action: get_entity_apis
Endpoint: /entities/{entity_id}/apis

CORRECT:
{
  "action": "get_entity_apis",
  "path_params": {
    "entity_id": "$2.result.id[]"
  }
}

WRONG:
{
  "query_params": {
    "entity_id": "$2.result.id[]"
  }
}

-------------------------

Return ONLY JSON:

{
  "plan": [
    {
      "step": 1,
      "action": "action_name",
      "query_params": {},
      "path_params": {},
      "body": {}
    }
  ]
}
`, query, string(actionsJSON))
}

func BuildSummaryPrompt(userQuery string, apiResults []string) string {
	return fmt.Sprintf(`
Instructions:
- Provide concise summary based on the following API call results only.
- Do NOT use any information outside of these results.
- If results are empty, say "No relevant information found".
- If the user hasn't specified format for data then use markdown tables for tabular data and JSON for non-tabular data.

User Query:
%s

Step Results:
%s

Provide final concise markdown summary.
`, userQuery, strings.Join(apiResults, "\n"))
}
