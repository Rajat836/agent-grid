package agent_config

const scopes = "entities:view entities:manage services:view features:view features:manage teams:view apis:view apis:manage"

var OntologyAgentActionsList = map[string]Action{
	"get_features": {
		Name:        "get_features",
		Title:       "Get Features",
		Description: "Retrieve a list of features",
		UserExamples: []string{
			"List all features",
			"Get active features",
		},
		Filters: []FilterParam{
			{Key: "code", Type: "string", Description: "Filter by feature code"},
			{Key: "is_active", Type: "bool", Description: "Filter active features"},
			{Key: "created_at_from", Type: "date", Description: "Start date"},
		},
		Pagination: true,
		ResponseJSON: `{
"response_type":"json",
"action":"get_features",
"filters":{"code":"<code>","is_active":"<true|false>"},
"pagination":{"page":1,"limit":10}
}`,
		API: APIConfig{
			Host:     "http://localhost:4441",
			Method:   "GET",
			Endpoint: "ontology/v1/features",
			Headers: map[string]string{
				"X-Scope": scopes,
			},
		},
	},
}
