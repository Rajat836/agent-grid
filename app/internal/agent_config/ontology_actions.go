package agent_config

import (
	"net/http"

	"bitbucket.org/fyscal/be-commons/pkg/global"
)

const scopes = "entities:view entities:manage services:view features:view features:manage teams:view apis:view apis:manage"

type ServiceHost string

const (
	OntologyServiceHost ServiceHost = "http://localhost:4441"
)

type ActionName string

const (
	// Features
	GetFeatures         ActionName = "get_features"
	CreateFeature       ActionName = "create_feature"
	UpdateFeature       ActionName = "update_feature"
	GetFeatureInstances ActionName = "get_feature_instances"
	GetFeatureMetrics   ActionName = "get_feature_metrics"

	// Entities
	GetEntities            ActionName = "get_entities"
	CreateEntity           ActionName = "create_entity"
	UpdateEntity           ActionName = "update_entity"
	GetEntityMetrics       ActionName = "get_entity_metrics"
	GetEntityAPIs          ActionName = "get_entity_apis"
	GetEntityTransitions   ActionName = "get_entity_transitions"
	CreateEntityTransition ActionName = "create_entity_transition"
	UpdateEntityTransition ActionName = "update_entity_transition"

	// Services
	GetServices           ActionName = "get_services"
	UpdateService         ActionName = "update_service"
	AssignServiceTeam     ActionName = "assign_service_team"
	GetServiceDeployments ActionName = "get_service_deployments"

	// Teams
	GetTeams          ActionName = "get_teams"
	CreateTeam        ActionName = "create_team"
	UpdateTeam        ActionName = "update_team"
	GetTeamsByFeature ActionName = "get_teams_by_feature"

	// APIs
	GetAPIs       ActionName = "get_apis"
	UpdateAPI     ActionName = "update_api"
	GetAPIMetrics ActionName = "get_api_metrics"

	// KPIs
	GetKPIs               ActionName = "get_kpis"
	CreateKPI             ActionName = "create_kpi"
	UpdateKPI             ActionName = "update_kpi"
	GetKPIRelationships   ActionName = "get_kpi_relationships"
	CreateKPIRelationship ActionName = "create_kpi_relationship"
	UpdateKPIRelationship ActionName = "update_kpi_relationship"
)

var OntologyAgentActionsList = map[ActionName]Action{

	// =========================
	// GET FEATURES
	// =========================
	GetFeatures: {
		Name:        GetFeatures,
		Title:       "Get Features",
		Description: "Fetch list of features with metadata like id, code, name, description, status and timestamps. Use this when user asks to list/search/filter features.",
		UserExamples: []string{
			"List all features",
			"Get active features",
			"Find feature with code xyz",
			"Show features created after Feb 2026",
		},

		QueryParams: []Param{
			{"code", "string", false, "Filter by feature code", "xyz"},
			{"name", "string", false, "Filter by feature name", "Customer Onboarding"},
			{"is_active", "bool", false, "Filter active/inactive", "true"},
			{"created_at_from", "date", false, "Start date filter", "2026-01-01"},
			{"created_at_to", "date", false, "End date filter", "2026-02-01"},
			{"page", "int", false, "Page number", "1"},
			{"limit", "int", false, "Page size", "10"},
		},

		Pagination: true,

		ResponseJSON: `{
"action":"get_features",
"query_params":{
  "code":"<code>",
  "name":"<name>",
  "is_active":"<true|false>"
}
}`,

		API: APIConfig{
			Host:     OntologyServiceHost,
			Method:   http.MethodGet,
			Endpoint: "ontology/v1/features",
			Headers: map[string]string{
				string(global.XScope): scopes,
			},
		},
	},

	// =========================
	// CREATE FEATURE
	// =========================
	CreateFeature: {
		Name:        CreateFeature,
		Title:       "Create Feature",
		Description: "Create a new feature with code, name, description and active status. Use when user wants to add/register a new feature.",
		UserExamples: []string{
			"Create a new feature for onboarding",
			"Add feature xyz",
			"Register a feature for user signup flow",
		},

		BodyParams: []Param{
			{"code", "string", true, "Unique feature code", "xyz"},
			{"name", "string", true, "Feature name", "Customer Onboarding"},
			{"description", "string", false, "Feature description", "Flow for onboarding"},
			{"is_active", "bool", true, "Feature status", "true"},
		},

		ResponseJSON: `{
"action":"create_feature",
"body":{
  "code":"xyz",
  "name":"Customer Onboarding"
}
}`,

		API: APIConfig{
			Host:     OntologyServiceHost,
			Method:   http.MethodPost,
			Endpoint: "ontology/v1/features",
			Headers: map[string]string{
				string(global.XScope): scopes,
			},
		},
	},

	// =========================
	// UPDATE FEATURE
	// =========================
	UpdateFeature: {
		Name:        UpdateFeature,
		Title:       "Update Feature",
		Description: "Update an existing feature using feature_id. Use when modifying name, description or active status.",
		UserExamples: []string{
			"Update feature 101 name",
			"Disable feature xyz",
			"Change description of onboarding feature",
		},

		PathParams: []Param{
			{"feature_id", "int", true, "Feature ID", "101"},
		},

		BodyParams: []Param{
			{"name", "string", false, "Feature name", "Customer Onboarding"},
			{"description", "string", false, "Feature description", "Updated flow"},
			{"is_active", "bool", false, "Active status", "true"},
		},

		ResponseJSON: `{
"action":"update_feature",
"path_params":{"feature_id":101},
"body":{"name":"Updated Name"}
}`,

		API: APIConfig{
			Host:     OntologyServiceHost,
			Method:   http.MethodPatch,
			Endpoint: "ontology/v1/features/{feature_id}",
			Headers: map[string]string{
				string(global.XScope): scopes,
			},
		},
	},

	// =========================
	// FEATURE INSTANCES
	// =========================
	GetFeatureInstances: {
		Name:        GetFeatureInstances,
		Title:       "Get Feature Instances",
		Description: "Fetch execution instances of features with status and timestamps. Use when user asks about feature runs, executions, or history.",
		UserExamples: []string{
			"Show feature executions",
			"Get completed feature runs",
			"List instances for feature 101",
		},

		QueryParams: []Param{
			{"feature_id", "int", false, "Filter by feature id", "101"},
			{"status", "string", false, "Execution status", "completed"},
			{"started_at_from", "date", false, "Start time filter", ""},
			{"completed_at_to", "date", false, "End time filter", ""},
			{"limit", "int", false, "Limit results", "100"},
		},

		ResponseJSON: `{
"action":"get_feature_instances",
"query_params":{"feature_id":101}
}`,

		API: APIConfig{
			Host:     OntologyServiceHost,
			Method:   http.MethodGet,
			Endpoint: "ontology/v1/features/instances",
			Headers: map[string]string{
				string(global.XScope): scopes,
			},
		},
	},

	// =========================
	// FEATURE METRICS
	// =========================
	GetFeatureMetrics: {
		Name:        GetFeatureMetrics,
		Title:       "Get Feature Metrics",
		Description: "Fetch performance metrics like success rate, latency (p50, p95, p99) for a feature. Use when user asks for analytics or performance.",
		UserExamples: []string{
			"Show metrics for feature 101",
			"Get success rate of onboarding feature",
			"Latency stats for feature xyz",
		},

		PathParams: []Param{
			{"feature_id", "int", true, "Feature ID", "101"},
		},

		QueryParams: []Param{
			{"window_start_from", "date", false, "Start window", ""},
			{"window_start_to", "date", false, "End window", ""},
			{"window_minutes", "int", false, "Aggregation window", "5"},
		},

		ResponseJSON: `{
"action":"get_feature_metrics",
"path_params":{"feature_id":101}
}`,

		API: APIConfig{
			Host:     OntologyServiceHost,
			Method:   http.MethodGet,
			Endpoint: "ontology/v1/features/{feature_id}/metrics",
			Headers: map[string]string{
				string(global.XScope): scopes,
			},
		},
	},

	//
	// =========================
	// ENTITIES
	// =========================
	//

	GetEntities: {
		Name:        GetEntities,
		Title:       "Get Entities",
		Description: "Fetch workflow entities (steps) within a feature including ordering, start/end flags and metadata. Use when user asks about steps, stages, or flow structure of a feature.",
		UserExamples: []string{
			"List entities for feature 101",
			"Show all steps in onboarding flow",
			"Get start entities",
			"Find terminal entities",
		},

		QueryParams: []Param{
			{"feature_id", "int", false, "Filter by feature id", "101"},
			{"code", "string", false, "Entity code", "USER_REG"},
			{"name", "string", false, "Entity name", "User Registration"},
			{"is_start", "bool", false, "Start node filter", "true"},
			{"is_terminal", "bool", false, "Terminal node filter", "false"},
			{"created_at_from", "date", false, "Start date", ""},
			{"created_at_to", "date", false, "End date", ""},
			{"page", "int", false, "Page number", "1"},
			{"limit", "int", false, "Page size", "10"},
		},

		Pagination: true,

		ResponseJSON: `{
"action":"get_entities",
"query_params":{"feature_id":101}
}`,

		API: APIConfig{
			Host:     OntologyServiceHost,
			Method:   http.MethodGet,
			Endpoint: "ontology/v1/entities",
			Headers: map[string]string{
				string(global.XScope): scopes,
			},
		},
	},

	// =========================
	// CREATE ENTITY
	// =========================
	CreateEntity: {
		Name:        CreateEntity,
		Title:       "Create Entity",
		Description: "Create a new workflow step (entity) inside a feature. Use when defining new steps in a flow.",
		UserExamples: []string{
			"Create a new step in onboarding",
			"Add entity USER_REG to feature 101",
			"Add a start node for signup flow",
		},

		BodyParams: []Param{
			{"feature_id", "int", true, "Feature ID", "101"},
			{"code", "string", true, "Entity code", "USER_REG"},
			{"name", "string", true, "Entity name", "User Registration"},
			{"description", "string", false, "Description", ""},
			{"display_order", "int", true, "Execution order", "1"},
			{"is_start", "bool", true, "Is start node", "true"},
			{"is_terminal", "bool", true, "Is terminal node", "false"},
		},

		ResponseJSON: `{
"action":"create_entity",
"body":{"feature_id":101,"code":"USER_REG"}
}`,

		API: APIConfig{
			Host:     OntologyServiceHost,
			Method:   http.MethodPost,
			Endpoint: "ontology/v1/entities",
			Headers: map[string]string{
				string(global.XScope): scopes,
			},
		},
	},

	// =========================
	// UPDATE ENTITY
	// =========================
	UpdateEntity: {
		Name:        UpdateEntity,
		Title:       "Update Entity",
		Description: "Update an existing entity's metadata like name, description, order or flags. Use when modifying workflow steps.",
		UserExamples: []string{
			"Update entity 501",
			"Change order of USER_REG step",
			"Mark entity as terminal",
		},

		PathParams: []Param{
			{"entity_id", "int", true, "Entity ID", "501"},
		},

		BodyParams: []Param{
			{"name", "string", false, "Entity name", ""},
			{"description", "string", false, "Description", ""},
			{"display_order", "int", false, "Order", "2"},
			{"is_start", "bool", false, "Start flag", "true"},
			{"is_terminal", "bool", false, "Terminal flag", "false"},
		},

		ResponseJSON: `{
"action":"update_entity",
"path_params":{"entity_id":501}
}`,

		API: APIConfig{
			Host:     OntologyServiceHost,
			Method:   http.MethodPatch,
			Endpoint: "ontology/v1/entities/{entity_id}",
			Headers: map[string]string{
				string(global.XScope): scopes,
			},
		},
	},

	// =========================
	// ENTITY METRICS
	// =========================
	GetEntityMetrics: {
		Name:        GetEntityMetrics,
		Title:       "Get Entity Metrics",
		Description: "Fetch execution metrics of an entity like success rate and latency. Use when analyzing performance of a workflow step.",
		UserExamples: []string{
			"Show metrics for entity 501",
			"Success rate of USER_REG step",
			"Latency of registration step",
		},

		PathParams: []Param{
			{"entity_id", "int", true, "Entity ID", "501"},
		},

		QueryParams: []Param{
			{"window_start_from", "date", false, "", ""},
			{"window_start_to", "date", false, "", ""},
			{"window_minutes", "int", false, "", "5"},
		},

		ResponseJSON: `{
"action":"get_entity_metrics",
"path_params":{"entity_id":501}
}`,

		API: APIConfig{
			Host:     OntologyServiceHost,
			Method:   http.MethodGet,
			Endpoint: "ontology/v1/entities/{entity_id}/metrics",
			Headers: map[string]string{
				string(global.XScope): scopes,
			},
		},
	},

	// =========================
	// ENTITY APIs
	// =========================
	GetEntityAPIs: {
		Name:        GetEntityAPIs,
		Title:       "Get Entity APIs",
		Description: "Fetch APIs associated with an entity including service, endpoint and method. Use when user asks which APIs are executed in a step.",
		UserExamples: []string{
			"What APIs are used in entity 501",
			"Show backend calls for USER_REG",
			"List APIs for registration step",
		},

		PathParams: []Param{
			{"entity_id", "int", true, "Entity ID", "501"},
		},

		Pagination: true,

		ResponseJSON: `{
"action":"get_entity_apis",
"path_params":{"entity_id":501}
}`,

		API: APIConfig{
			Host:     OntologyServiceHost,
			Method:   http.MethodGet,
			Endpoint: "ontology/v1/entities/{entity_id}/apis",
			Headers: map[string]string{
				string(global.XScope): scopes,
			},
		},
	},

	// =========================
	// ENTITY TRANSITIONS
	// =========================
	GetEntityTransitions: {
		Name:        GetEntityTransitions,
		Title:       "Get Entity Transitions",
		Description: "Fetch transitions between entities in a workflow including success/failure paths. Use when user asks about flow movement or state transitions.",
		UserExamples: []string{
			"Show transitions for feature 101",
			"How does USER_REG move to next step",
			"Get failure transitions",
		},

		QueryParams: []Param{
			{"feature_id", "int", false, "", "101"},
			{"from_entity_id", "int", false, "", ""},
			{"to_entity_id", "int", false, "", ""},
			{"transition_type", "string", false, "success/failure", "success"},
		},

		ResponseJSON: `{
"action":"get_entity_transitions",
"query_params":{"feature_id":101}
}`,

		API: APIConfig{
			Host:     OntologyServiceHost,
			Method:   http.MethodGet,
			Endpoint: "ontology/v1/entities/transitions",
			Headers: map[string]string{
				string(global.XScope): scopes,
			},
		},
	},

	// =========================
	// CREATE TRANSITION
	// =========================
	CreateEntityTransition: {
		Name:        CreateEntityTransition,
		Title:       "Create Entity Transition",
		Description: "Create a transition between two entities defining workflow movement.",
		UserExamples: []string{
			"Connect USER_REG to VERIFY_OTP",
			"Add success transition between entities",
		},

		BodyParams: []Param{
			{"feature_id", "int", true, "", "101"},
			{"from_entity_id", "int", true, "", "501"},
			{"to_entity_id", "int", true, "", "502"},
			{"transition_type", "string", true, "success/failure", "success"},
			{"condition_description", "string", false, "", ""},
		},

		ResponseJSON: `{
"action":"create_entity_transition"
}`,

		API: APIConfig{
			Host:     OntologyServiceHost,
			Method:   http.MethodPost,
			Endpoint: "ontology/v1/entities/transitions",
			Headers: map[string]string{
				string(global.XScope): scopes,
			},
		},
	},

	// =========================
	// UPDATE TRANSITION
	// =========================
	UpdateEntityTransition: {
		Name:        UpdateEntityTransition,
		Title:       "Update Entity Transition",
		Description: "Update an existing transition between entities.",
		UserExamples: []string{
			"Change transition 1002 to failure",
			"Update condition for entity transition",
		},

		PathParams: []Param{
			{"transition_id", "int", true, "", "1002"},
		},

		BodyParams: []Param{
			{"from_entity_id", "int", false, "", ""},
			{"to_entity_id", "int", false, "", ""},
			{"transition_type", "string", false, "", "failure"},
			{"condition_description", "string", false, "", ""},
		},

		ResponseJSON: `{
"action":"update_entity_transition",
"path_params":{"transition_id":1002}
}`,

		API: APIConfig{
			Host:     OntologyServiceHost,
			Method:   http.MethodPatch,
			Endpoint: "ontology/v1/entities/transitions/{transition_id}",
			Headers: map[string]string{
				string(global.XScope): scopes,
			},
		},
	},

	//
	// =========================
	// SERVICES
	// =========================
	//

	// =========================
	// GET SERVICES
	// =========================
	GetServices: {
		Name:        GetServices,
		Title:       "Get Services",
		Description: "Fetch backend services including code, name, criticality level, and repository details. Use when user asks about services, ownership, or repositories.",
		UserExamples: []string{
			"List all services",
			"Show critical services",
			"Find service AUTH_SVC",
			"Get services with repo github.com/company/auth-service",
		},

		QueryParams: []Param{
			{"code", "string", false, "Service code", "AUTH_SVC"},
			{"name", "string", false, "Service name", "Auth Service"},
			{"criticality_level", "string", false, "LOW | MEDIUM | HIGH", "HIGH"},
			{"repository_url", "string", false, "Repository URL", "github.com/company/auth-service"},
			{"created_at_from", "date", false, "Start date", ""},
			{"created_at_to", "date", false, "End date", ""},
			{"page", "int", false, "Page number", "1"},
			{"limit", "int", false, "Page size", "10"},
		},

		Pagination: true,

		ResponseJSON: `{
"action":"get_services",
"query_params":{"code":"AUTH_SVC"}
}`,

		API: APIConfig{
			Host:     OntologyServiceHost,
			Method:   http.MethodGet,
			Endpoint: "ontology/v1/services",
			Headers: map[string]string{
				string(global.XScope): scopes,
			},
		},
	},

	// =========================
	// UPDATE SERVICE
	// =========================
	UpdateService: {
		Name:        UpdateService,
		Title:       "Update Service",
		Description: "Update service metadata such as name, description, repository URL, and criticality level. Use when modifying service details.",
		UserExamples: []string{
			"Update service 1",
			"Change repo of AUTH_SVC",
			"Update criticality of auth service to HIGH",
		},

		PathParams: []Param{
			{"service_id", "int", true, "Service ID", "1"},
		},

		BodyParams: []Param{
			{"name", "string", false, "Service name", ""},
			{"description", "string", false, "Service description", ""},
			{"repository_url", "string", false, "Repository URL", ""},
			{"criticality_level", "string", false, "LOW | MEDIUM | HIGH", "HIGH"},
		},

		ResponseJSON: `{
"action":"update_service",
"path_params":{"service_id":1}
}`,

		API: APIConfig{
			Host:     OntologyServiceHost,
			Method:   http.MethodPatch,
			Endpoint: "ontology/v1/services/{service_id}",
			Headers: map[string]string{
				string(global.XScope): scopes,
			},
		},
	},

	// =========================
	// ASSIGN TEAM TO SERVICE
	// =========================
	AssignServiceTeam: {
		Name:        AssignServiceTeam,
		Title:       "Assign Team to Service",
		Description: "Assign a team to a service with a role such as owner. Use when defining ownership or responsibility.",
		UserExamples: []string{
			"Assign team 10 as owner to service 1",
			"Who owns auth service",
			"Add team to AUTH_SVC",
		},

		PathParams: []Param{
			{"service_id", "int", true, "Service ID", "1"},
		},

		BodyParams: []Param{
			{"team_id", "int", true, "Team ID", "10"},
			{"role", "string", true, "Role (owner/member)", "owner"},
		},

		ResponseJSON: `{
"action":"assign_service_team",
"path_params":{"service_id":1},
"body":{"team_id":10,"role":"owner"}
}`,

		API: APIConfig{
			Host:     OntologyServiceHost,
			Method:   http.MethodPost,
			Endpoint: "ontology/v1/services/{service_id}/teams",
			Headers: map[string]string{
				string(global.XScope): scopes,
			},
		},
	},

	// =========================
	// SERVICE DEPLOYMENTS
	// =========================
	GetServiceDeployments: {
		Name:        GetServiceDeployments,
		Title:       "Get Service Deployments",
		Description: "Fetch deployment history of services including environment, version and commit details. Use when user asks about releases, versions, or deployments.",
		UserExamples: []string{
			"Show deployments of auth service",
			"Latest prod deployment",
			"Which version is deployed in prod",
			"Get deployments for service 1",
		},

		QueryParams: []Param{
			{"service_id", "int", false, "Service ID", "1"},
			{"environment", "string", false, "Environment (dev/staging/prod)", "prod"},
			{"version", "string", false, "Version", "1.0.2"},
			{"commit_hash", "string", false, "Commit hash", "abc123"},
			{"deployed_at_from", "date", false, "Start date", ""},
			{"deployed_at_to", "date", false, "End date", ""},
			{"limit", "int", false, "Limit results", "100"},
		},

		ResponseJSON: `{
"action":"get_service_deployments",
"query_params":{"service_id":1}
}`,

		API: APIConfig{
			Host:     OntologyServiceHost,
			Method:   http.MethodGet,
			Endpoint: "ontology/v1/services/deployments",
			Headers: map[string]string{
				string(global.XScope): scopes,
			},
		},
	},

	//
	// =========================
	// TEAMS
	// =========================
	//

	// GET TEAMS
	GetTeams: {
		Name:        GetTeams,
		Title:       "Get Teams",
		Description: "Fetch teams responsible for services or features including slack channel, oncall email and pagerduty. Use when user asks about ownership or team details.",
		UserExamples: []string{
			"List all teams",
			"Who is oncall for platform",
			"Find team with slack #platform",
		},

		QueryParams: []Param{
			{"name", "string", false, "Team name", "Platform Team"},
			{"slack_channel", "string", false, "Slack channel", "#platform"},
			{"oncall_email", "string", false, "Oncall email", "platform@company.com"},
			{"created_at_from", "date", false, "", ""},
			{"created_at_to", "date", false, "", ""},
			{"page", "int", false, "Page number", "1"},
			{"limit", "int", false, "Page size", "10"},
		},

		Pagination: true,

		ResponseJSON: `{
"action":"get_teams",
"query_params":{"name":"Platform Team"}
}`,

		API: APIConfig{
			Host:     OntologyServiceHost,
			Method:   http.MethodGet,
			Endpoint: "ontology/v1/teams",
			Headers: map[string]string{
				string(global.XScope): scopes,
			},
		},
	},

	// CREATE TEAM
	CreateTeam: {
		Name:        CreateTeam,
		Title:       "Create Team",
		Description: "Create a new team with communication and oncall details.",
		UserExamples: []string{
			"Create platform team",
			"Add new team for backend",
		},

		BodyParams: []Param{
			{"name", "string", true, "", "Platform Team"},
			{"slack_channel", "string", true, "", "#platform"},
			{"pagerduty_key", "string", false, "", ""},
			{"oncall_email", "string", true, "", "platform@company.com"},
		},

		ResponseJSON: `{
"action":"create_team"
}`,

		API: APIConfig{
			Host:     OntologyServiceHost,
			Method:   http.MethodPost,
			Endpoint: "ontology/v1/teams",
			Headers: map[string]string{
				string(global.XScope): scopes,
			},
		},
	},

	// UPDATE TEAM
	UpdateTeam: {
		Name:        UpdateTeam,
		Title:       "Update Team",
		Description: "Update team metadata such as slack channel or oncall email.",
		UserExamples: []string{
			"Update team 10",
			"Change oncall email of platform team",
		},

		PathParams: []Param{
			{"team_id", "int", true, "", "10"},
		},

		BodyParams: []Param{
			{"name", "string", false, "", ""},
			{"slack_channel", "string", false, "", ""},
			{"pagerduty_key", "string", false, "", ""},
			{"oncall_email", "string", false, "", ""},
		},

		ResponseJSON: `{
"action":"update_team",
"path_params":{"team_id":10}
}`,

		API: APIConfig{
			Host:     OntologyServiceHost,
			Method:   http.MethodPatch,
			Endpoint: "ontology/v1/teams/{team_id}",
			Headers: map[string]string{
				string(global.XScope): scopes,
			},
		},
	},

	// TEAMS BY FEATURE
	GetTeamsByFeature: {
		Name:        GetTeamsByFeature,
		Title:       "Get Teams by Feature",
		Description: "Fetch teams responsible for a feature including their roles. Use when user asks ownership of a feature.",
		UserExamples: []string{
			"Who owns feature 101",
			"Teams responsible for onboarding",
		},

		PathParams: []Param{
			{"feature_id", "int", true, "", "101"},
		},

		ResponseJSON: `{
"action":"get_teams_by_feature",
"path_params":{"feature_id":101}
}`,

		API: APIConfig{
			Host:     OntologyServiceHost,
			Method:   http.MethodGet,
			Endpoint: "ontology/v1/teams/feature/{feature_id}",
			Headers: map[string]string{
				string(global.XScope): scopes,
			},
		},
	},

	//
	// =========================
	// APIs
	// =========================
	//

	// GET APIS
	GetAPIs: {
		Name:        GetAPIs,
		Title:       "Get APIs",
		Description: "Fetch API definitions including endpoint, method and protocol. Use when user asks about API catalog or endpoints.",
		UserExamples: []string{
			"List APIs",
			"Find POST APIs",
			"Show external APIs",
		},

		QueryParams: []Param{
			{"id", "int", false, "", "9001"},
			{"endpoint", "string", false, "", "/api/v1/users/register"},
			{"method", "string", false, "", "POST"},
			{"protocol", "string", false, "", "HTTP"},
			{"is_internal", "bool", false, "", "false"},
			{"created_at_from", "date", false, "", ""},
			{"created_at_to", "date", false, "", ""},
			{"page", "int", false, "", "1"},
			{"limit", "int", false, "", "10"},
		},

		Pagination: true,

		ResponseJSON: `{
"action":"get_apis",
"query_params":{"method":"POST"}
}`,

		API: APIConfig{
			Host:     OntologyServiceHost,
			Method:   http.MethodGet,
			Endpoint: "ontology/v1/apis",
			Headers: map[string]string{
				string(global.XScope): scopes,
			},
		},
	},

	// UPDATE API
	UpdateAPI: {
		Name:        UpdateAPI,
		Title:       "Update API",
		Description: "Update API metadata such as description.",
		UserExamples: []string{
			"Update API 9001 description",
			"Describe register API",
		},

		PathParams: []Param{
			{"api_id", "int", true, "", "9001"},
		},

		BodyParams: []Param{
			{"description", "string", true, "", "Registers user"},
		},

		ResponseJSON: `{
"action":"update_api",
"path_params":{"api_id":9001}
}`,

		API: APIConfig{
			Host:     OntologyServiceHost,
			Method:   http.MethodPatch,
			Endpoint: "ontology/v1/apis/{api_id}",
			Headers: map[string]string{
				string(global.XScope): scopes,
			},
		},
	},

	// API METRICS
	GetAPIMetrics: {
		Name:        GetAPIMetrics,
		Title:       "Get API Metrics",
		Description: "Fetch performance metrics for an API such as latency and success rate.",
		UserExamples: []string{
			"Show metrics for API 9001",
			"Latency of register API",
		},

		PathParams: []Param{
			{"api_id", "int", true, "", "9001"},
		},

		QueryParams: []Param{
			{"apis_deployment_id", "int", false, "", ""},
			{"window_start_from", "date", false, "", ""},
			{"window_start_to", "date", false, "", ""},
			{"window_minutes", "int", false, "", "5"},
		},

		ResponseJSON: `{
"action":"get_api_metrics",
"path_params":{"api_id":9001}
}`,

		API: APIConfig{
			Host:     OntologyServiceHost,
			Method:   http.MethodGet,
			Endpoint: "ontology/v1/apis/{api_id}/metrics",
			Headers: map[string]string{
				string(global.XScope): scopes,
			},
		},
	},

	//
	// =========================
	// KPIs
	// =========================
	//

	// GET KPIs
	GetKPIs: {
		Name:        GetKPIs,
		Title:       "Get KPIs",
		Description: "Fetch KPI definitions such as latency or success rate metrics.",
		UserExamples: []string{
			"List KPIs",
			"Show latency metrics",
		},

		QueryParams: []Param{
			{"code", "string", false, "", "API_LATENCY"},
			{"name", "string", false, "", ""},
			{"metric_type", "string", false, "", "latency"},
			{"unit", "string", false, "", "ms"},
			{"created_at_from", "date", false, "", ""},
			{"created_at_to", "date", false, "", ""},
			{"page", "int", false, "", "1"},
			{"limit", "int", false, "", "10"},
		},

		Pagination: true,

		ResponseJSON: `{
"action":"get_kpis"
}`,

		API: APIConfig{
			Host:     OntologyServiceHost,
			Method:   http.MethodGet,
			Endpoint: "ontology/v1/kpis",
			Headers: map[string]string{
				string(global.XScope): scopes,
			},
		},
	},

	// CREATE KPI
	CreateKPI: {
		Name:        CreateKPI,
		Title:       "Create KPI",
		Description: "Create a new KPI definition.",
		UserExamples: []string{
			"Create API latency KPI",
		},

		BodyParams: []Param{
			{"code", "string", true, "", "API_LATENCY"},
			{"name", "string", true, "", "API Latency"},
			{"description", "string", false, "", ""},
			{"metric_type", "string", true, "", "latency"},
			{"unit", "string", true, "", "ms"},
		},

		ResponseJSON: `{
"action":"create_kpi"
}`,

		API: APIConfig{
			Host:     OntologyServiceHost,
			Method:   http.MethodPost,
			Endpoint: "ontology/v1/kpis",
			Headers: map[string]string{
				string(global.XScope): scopes,
			},
		},
	},

	// UPDATE KPI
	UpdateKPI: {
		Name:        UpdateKPI,
		Title:       "Update KPI",
		Description: "Update KPI metadata.",
		UserExamples: []string{
			"Update KPI 1",
		},

		PathParams: []Param{
			{"kpi_id", "int", true, "", "1"},
		},

		BodyParams: []Param{
			{"name", "string", false, "", ""},
			{"description", "string", false, "", ""},
			{"metric_type", "string", false, "", ""},
			{"unit", "string", false, "", ""},
		},

		ResponseJSON: `{
"action":"update_kpi",
"path_params":{"kpi_id":1}
}`,

		API: APIConfig{
			Host:     OntologyServiceHost,
			Method:   http.MethodPatch,
			Endpoint: "ontology/v1/kpis/{kpi_id}",
			Headers: map[string]string{
				string(global.XScope): scopes,
			},
		},
	},

	// KPI RELATIONSHIPS
	GetKPIRelationships: {
		Name:        GetKPIRelationships,
		Title:       "Get KPI Relationships",
		Description: "Fetch relationships showing how KPIs impact features, entities or services.",
		UserExamples: []string{
			"Show KPI relationships",
			"Which KPI impacts feature 101",
		},

		QueryParams: []Param{
			{"kpi_id", "int", false, "", "1"},
			{"relation_type", "string", false, "", "impacts"},
			{"target_type", "string", false, "", "feature"},
			{"target_id", "int", false, "", "101"},
		},

		ResponseJSON: `{
"action":"get_kpi_relationships"
}`,

		API: APIConfig{
			Host:     OntologyServiceHost,
			Method:   http.MethodGet,
			Endpoint: "ontology/v1/kpi-relationships",
			Headers: map[string]string{
				string(global.XScope): scopes,
			},
		},
	},

	// CREATE KPI RELATIONSHIP
	CreateKPIRelationship: {
		Name:        CreateKPIRelationship,
		Title:       "Create KPI Relationship",
		Description: "Create a relationship linking KPI to feature/entity/service.",
		UserExamples: []string{
			"Link API latency to feature 101",
		},

		BodyParams: []Param{
			{"kpi_id", "int", true, "", "1"},
			{"relation_type", "string", true, "", "impacts"},
			{"target_type", "string", true, "", "feature"},
			{"target_id", "int", true, "", "101"},
			{"weight", "float", true, "", "0.4"},
			{"weight_set_by", "string", true, "", "rajat@fyscaltech.com"},
		},

		ResponseJSON: `{
"action":"create_kpi_relationship"
}`,

		API: APIConfig{
			Host:     OntologyServiceHost,
			Method:   http.MethodPost,
			Endpoint: "ontology/v1/kpi-relationships",
			Headers: map[string]string{
				string(global.XScope): scopes,
			},
		},
	},

	// UPDATE KPI RELATIONSHIP
	UpdateKPIRelationship: {
		Name:        UpdateKPIRelationship,
		Title:       "Update KPI Relationship",
		Description: "Update KPI relationship weights or mapping.",
		UserExamples: []string{
			"Update KPI relationship 10",
		},

		PathParams: []Param{
			{"relationship_id", "int", true, "", "10"},
		},

		BodyParams: []Param{
			{"relation_type", "string", false, "", ""},
			{"target_type", "string", false, "", ""},
			{"target_id", "int", false, "", ""},
			{"weight", "float", false, "", "0.6"},
			{"weight_set_by", "string", false, "", ""},
			{"weight_reviewed_by", "string", false, "", ""},
		},

		ResponseJSON: `{
"action":"update_kpi_relationship",
"path_params":{"relationship_id":10}
}`,

		API: APIConfig{
			Host:     OntologyServiceHost,
			Method:   http.MethodPatch,
			Endpoint: "ontology/v1/kpi-relationships/{relationship_id}",
			Headers: map[string]string{
				string(global.XScope): scopes,
			},
		},
	},
}
