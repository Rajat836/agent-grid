package agent_config

import (
	"encoding/json"
	"fmt"
	"strings"
)

// =========================
// LLM STRUCTS
// =========================

// type LLMAction struct {
// 	Name         string   `json:"name"`
// 	Description  string   `json:"description"`
// 	Params       []string `json:"params"`
// 	UserExamples []string `json:"user_examples,omitempty"`
// 	RequestBody  any      `json:"request_body,omitempty"`
// }

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

var DomainActionsMap = map[string][]ActionName{
	"features": {
		GetFeatures,
		// CreateFeature,
		// UpdateFeature,
		// GetFeatureInstances,
		// GetFeatureMetrics,
	},
	"entities": {
		GetEntities,
		CreateEntity,
		UpdateEntity,
		GetEntityMetrics,
		GetEntityAPIs,
		// GetEntityTransitions,
		// CreateEntityTransition,
		// UpdateEntityTransition,
	},
	"services": {
		GetServices,
		// UpdateService,
		// AssignServiceTeam,
		// GetServiceDeployments,
	},
	"teams": {
		GetTeams,
		// CreateTeam,
		// UpdateTeam,
		// GetTeamsByFeature,
	},
	"apis": {
		GetAPIs,
		// UpdateAPI,
		GetAPIMetrics,
	},
	"kpis": {
		GetKPIs,
		// CreateKPI,
		// UpdateKPI,
		GetKPIRelationships,
		// CreateKPIRelationship,
		// UpdateKPIRelationship,
	},
}

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
	return DomainActionsMap
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
	// domains := DetectDomains(query)
	domainActionMap := createDomainActionMap()

	actionMap := map[ActionName]Action{}

	for _, v := range domainActionMap {
		for _, a := range v {
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
type LLMParam struct {
	Key         string `json:"key"`
	Type        string `json:"type"`
	Required    bool   `json:"required"`
	Description string `json:"description"`
	Example     any    `json:"example"`
}

type LLMAction struct {
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	UserExamples []string   `json:"user_examples,omitempty"`
	PathParams   []LLMParam `json:"path_params,omitempty"`
	QueryParams  []LLMParam `json:"query_params,omitempty"`
	BodyParams   []LLMParam `json:"body_params,omitempty"`
}

func ToLLMActions(actions []Action) []LLMAction {
	var result []LLMAction

	for _, a := range actions {
		llm := LLMAction{
			Name:         string(a.Name),
			Description:  a.Description,
			UserExamples: a.UserExamples,
		}

		for _, p := range a.PathParams {
			llm.PathParams = append(llm.PathParams, LLMParam{
				Key:         p.Key,
				Type:        p.Type,
				Required:    p.Required,
				Description: p.Description,
				Example:     p.Example,
			})
		}

		for _, p := range a.QueryParams {
			llm.QueryParams = append(llm.QueryParams, LLMParam{
				Key:         p.Key,
				Type:        p.Type,
				Required:    p.Required,
				Description: p.Description,
				Example:     p.Example,
			})
		}

		for _, p := range a.BodyParams {
			llm.BodyParams = append(llm.BodyParams, LLMParam{
				Key:         p.Key,
				Type:        p.Type,
				Required:    p.Required,
				Description: p.Description,
				Example:     p.Example,
			})
		}

		result = append(result, llm)
	}

	return result
}

// =========================
// PLANNER PROMPT
// =========================
func BuildPlannerPrompt(query string, actions []LLMAction) string {
	actionsJSON, _ := json.MarshalIndent(actions, "", "  ")

	return fmt.Sprintf(`
You are an API orchestration planner. Given a user query and a list of available actions, return a strict JSON execution plan.

=========================
USER QUERY
=========================
%s

=========================
AVAILABLE ACTIONS
=========================
Each action defines exactly which params belong in path_params, query_params, or body_params.
You MUST place each param in the bucket it is listed under. Never move a param to a different bucket.

%s

=========================
OUTPUT FORMAT
=========================

Return ONLY this JSON — no explanation, no markdown, no extra text:

{
  "plan": [
    {
      "step": 1,
      "action": "<action_name>",
      "reason": "<why this step is needed>",
      "depends_on": null,
      "endpoint": "<endpoint path for logging>",
      "query_params": {},
      "path_params": {},
      "body": {}
    }
  ]
}

=========================
STRICT RULES
=========================

RULE 1 — Only use actions from the provided list. Never invent action names.

RULE 2 — Every step MUST have all fields: step, action, reason, depends_on, endpoint, query_params, path_params, body.
  - Use {} for empty param objects.
  - Use null for depends_on when there is no dependency.

RULE 3 — depends_on:
  - Set to the step number this step must wait for.
  - null if the step has no dependency.
  - A step can only depend on ONE prior step.

RULE 4 — Required params:
  - Any param marked "required": true MUST be present in the step.
  - Never omit a required param.

=========================
PARAM PLACEMENT RULES
=========================

query_params:
  - Only params listed under "query_params" in the action definition.
  - These become URL filters: ?key=value
  - Example: { "code": "customer_onboarding", "is_active": true }

path_params:
  - Only params listed under "path_params" in the action definition.
  - These fill {placeholder} variables in the endpoint path.
  - NEVER place path params inside query_params or body.
  - Example:
      Action path_params: [{ "key": "entity_id" }]
      Endpoint:           /entities/{entity_id}/apis
      Correct:            "path_params": { "entity_id": 42 }
      WRONG:              "query_params": { "entity_id": 42 }

body:
  - Only params listed under "body_params" in the action definition.
  - Only used for POST or PATCH actions.
  - Always {} for GET actions.

=========================
CHAINING — REFERENCING PRIOR STEPS
=========================

When a param value comes from a prior step's response, use this reference syntax:

  Single value:  "$<step>.result.<field>"
  Array of IDs:  "$<step>.result.<field>[]"

Use [] when the prior step returns multiple records and you need all their IDs.
Use without [] when the prior step returns a single record or one specific value.

Examples:
  "$1.result.id"         — single ID from step 1
  "$2.result.id[]"       — all IDs from step 2 (multiple records)
  "$1.result.feature_id" — a specific field from step 1

IMPORTANT: Chaining references follow the same placement rules.
  A referenced value goes in the SAME bucket its param key is listed under in the action definition.

=========================
MULTI-STEP CHAINING EXAMPLE
=========================

User query: "List all APIs in the customer_onboarding feature"

{
  "plan": [
    {
      "step": 1,
      "action": "get_features",
      "reason": "Resolve feature ID for code customer_onboarding",
      "depends_on": null,
      "endpoint": "/features",
      "query_params": { "code": "customer_onboarding" },
      "path_params": {},
      "body": {}
    },
    {
      "step": 2,
      "action": "get_entities",
      "reason": "Fetch all entities belonging to the resolved feature",
      "depends_on": 1,
      "endpoint": "/entities",
      "query_params": { "feature_id": "$1.result.id" },
      "path_params": {},
      "body": {}
    },
    {
      "step": 3,
      "action": "get_entity_apis",
      "reason": "Fetch APIs for every entity returned in step 2",
      "depends_on": 2,
      "endpoint": "/entities/{entity_id}/apis",
      "query_params": {},
      "path_params": { "entity_id": "$2.result.id[]" },
      "body": {}
    }
  ]
}

=========================
REMINDERS
=========================

- Place every param in the bucket its action definition lists it under.
- Never fabricate action names.
- Always set depends_on (null or a step number).
- Always write a reason for each step.
- Return ONLY the JSON object. No prose, no markdown fences.
`, query, string(actionsJSON))
}

func BuildSummaryPrompt(userQuery string, apiResults []string) string {
	return fmt.Sprintf(`
Instructions:
- Provide concise summary based on the following API call results only.
- Do NOT use any information outside of these results.
- If results are empty, say "No relevant information found".
- Use all the step results to provide comprehensive summary. Do NOT ignore any field n.
- If the user hasn't specified format for data then use markdown tables for tabular data and JSON for non-tabular data.

User Query:
%s

Step Results:
%s

Provide final concise markdown summary.
`, userQuery, strings.Join(apiResults, "\n"))
}
