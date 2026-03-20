package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"app/agent_grid/internal/config"

	"github.com/google/uuid"
	"gorm.io/driver/clickhouse"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// --- Struct Definitions ---

type ApiDbInteraction struct {
	ID              string  `gorm:"column:id" json:"id"`
	ApiExecutionID  string  `gorm:"column:api_execution_id" json:"api_execution_id"`
	DatabaseType    string  `gorm:"column:database_type;type:LowCardinality(String)" json:"database_type"`
	DbTableName     string  `gorm:"column:table_name" json:"table_name"`
	OperationType   string  `gorm:"column:operation_type;type:LowCardinality(String)" json:"operation_type"`
	SqlQuery        string  `gorm:"column:sql_query" json:"sql_query"`
	QueryParameters *string `gorm:"column:query_parameters" json:"query_parameters,omitempty"`
	RowsAffected    int32   `gorm:"column:rows_affected" json:"rows_affected"`
	DurationMs      int32   `gorm:"column:duration_ms" json:"duration_ms"`
	ErrorMessage    *string `gorm:"column:error_message" json:"error_message,omitempty"`
	CreatedAt       int64   `gorm:"column:created_at" json:"created_at"`
}

func (ApiDbInteraction) TableName() string { return "api_db_interactions" }

type ApiExecution struct {
	ID                string   `gorm:"column:id" json:"id"`
	EntityInstanceID  string   `gorm:"column:entity_instance_id" json:"entity_instance_id"`
	ApiDeploymentsID  string   `gorm:"column:api_deployments_id" json:"api_deployments_id"`
	OauthClientID     *string  `gorm:"column:oauth_client_id" json:"oauth_client_id,omitempty"`
	UserID            *string  `gorm:"column:user_id" json:"user_id,omitempty"`
	UserUID           *string  `gorm:"column:user_uid" json:"user_uid,omitempty"`
	OrganizationID    *string  `gorm:"column:organization_id" json:"organization_id,omitempty"`
	TokenID           *string  `gorm:"column:token_id" json:"token_id,omitempty"`
	RequestTimestamp  int64    `gorm:"column:request_timestamp" json:"request_timestamp"`
	RequestDurationMs int32    `gorm:"column:request_duration_ms" json:"request_duration_ms"`
	SourceIP          *string  `gorm:"column:source_ip" json:"source_ip,omitempty"`
	RequestID         *string  `gorm:"column:request_id" json:"request_id,omitempty"`
	CorrelationID     *string  `gorm:"column:correlation_id" json:"correlation_id,omitempty"`
	TraceID           *string  `gorm:"column:trace_id" json:"trace_id,omitempty"`
	SpanID            *string  `gorm:"column:span_id" json:"span_id,omitempty"`
	Host              *string  `gorm:"column:host" json:"host,omitempty"`
	Latitude          *float64 `gorm:"column:latitude" json:"latitude,omitempty"`
	Longitude         *float64 `gorm:"column:longitude" json:"longitude,omitempty"`
	DeviceID          *string  `gorm:"column:device_id" json:"device_id,omitempty"`
	PlatformOS        *string  `gorm:"column:platform_os;type:LowCardinality(Nullable(String))" json:"platform_os,omitempty"`
	AppVersion        *string  `gorm:"column:app_version" json:"app_version,omitempty"`
	UserAgent         *string  `gorm:"column:user_agent" json:"user_agent,omitempty"`
	RiskSessionID     *string  `gorm:"column:risk_session_id" json:"risk_session_id,omitempty"`
	DeviceModel       *string  `gorm:"column:device_model" json:"device_model,omitempty"`
	AuthScheme        *string  `gorm:"column:auth_scheme;type:LowCardinality(Nullable(String))" json:"auth_scheme,omitempty"`
	TokenType         *string  `gorm:"column:token_type;type:LowCardinality(Nullable(String))" json:"token_type,omitempty"`
	TokenScopes       *string  `gorm:"column:token_scopes" json:"token_scopes,omitempty"`
	HttpStatus        int32    `gorm:"column:http_status" json:"http_status"`
	AppErrorCode      *string  `gorm:"column:app_error_code" json:"app_error_code,omitempty"`
	ErrorMessage      *string  `gorm:"column:error_message" json:"error_message,omitempty"`
	ResponseSizeBytes int32    `gorm:"column:response_size_bytes" json:"response_size_bytes"`
	CreatedAt         int64    `gorm:"column:created_at" json:"created_at"`
}

func (ApiExecution) TableName() string { return "api_executions" }

type ApiMetric struct {
	ApisDeploymentID string    `gorm:"column:apis_deployment_id" json:"apis_deployment_id"`
	WindowStart      time.Time `gorm:"column:window_start" json:"window_start"`
	WindowMinutes    int32     `gorm:"column:window_minutes" json:"window_minutes"`
	TotalCalls       int64     `gorm:"column:total_calls" json:"total_calls"`
	SuccessRate      float64   `gorm:"column:success_rate" json:"success_rate"`
	ErrorRate        float64   `gorm:"column:error_rate" json:"error_rate"`
	P50LatencyMs     float64   `gorm:"column:p50_latency_ms" json:"p50_latency_ms"`
	P95LatencyMs     float64   `gorm:"column:p95_latency_ms" json:"p95_latency_ms"`
	P99LatencyMs     float64   `gorm:"column:p99_latency_ms" json:"p99_latency_ms"`
	CreatedAt        int64     `gorm:"column:created_at" json:"created_at"`
	UpdatedAt        int64     `gorm:"column:updated_at" json:"updated_at"`
}

func (ApiMetric) TableName() string { return "api_metrics" }

type Api struct {
	ID          int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Endpoint    string         `gorm:"type:varchar(512);not null" json:"endpoint"`
	HttpMethod  string         `gorm:"type:varchar(16);not null" json:"http_method"`
	Protocol    string         `gorm:"type:varchar(16);default:'HTTP'" json:"protocol"`
	IsInternal  bool           `gorm:"default:false;index:idx_apis_is_internal" json:"is_internal"`
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"default:NULL" json:"deleted_at"`
}

func (Api) TableName() string { return "apis" }

type ApisDeployment struct {
	ID                  string `gorm:"column:id" json:"id"`
	ApiID               int64  `gorm:"column:api_id" json:"api_id"`
	ServiceDeploymentID int64  `gorm:"column:service_deployment_id" json:"service_deployment_id"`
}

func (ApisDeployment) TableName() string { return "apis_deployments" }

type EntityApi struct {
	ID       int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	EntityID int64  `gorm:"not null;index:idx_entity_apis_entity_id" json:"entity_id"`
	ApiID    int64  `gorm:"not null;index:idx_entity_apis_api_id" json:"api_id"`
	Entity   Entity `gorm:"foreignKey:EntityID" json:"entity,omitempty"`
	Api      Api    `gorm:"foreignKey:ApiID" json:"api,omitempty"`
}

func (EntityApi) TableName() string { return "entity_apis" }

type EntityInstance struct {
	ID                string `gorm:"column:id" json:"id"`
	FeatureInstanceID string `gorm:"column:feature_instance_id" json:"feature_instance_id"`
	EntityID          int64  `gorm:"column:entity_id" json:"entity_id"`
	Status            string `gorm:"column:status;type:LowCardinality(String)" json:"status"`
	StartedAt         int64  `gorm:"column:started_at" json:"started_at"`
	CompletedAt       *int64 `gorm:"column:completed_at" json:"completed_at,omitempty"`
	CreatedAt         int64  `gorm:"column:created_at" json:"created_at"`
}

func (EntityInstance) TableName() string { return "entity_instances" }

type EntityMetric struct {
	EntityID      int64     `gorm:"column:entity_id" json:"entity_id"`
	WindowStart   time.Time `gorm:"column:window_start" json:"window_start"`
	WindowMinutes int32     `gorm:"column:window_minutes" json:"window_minutes"`
	TotalCount    int64     `gorm:"column:total_count" json:"total_count"`
	SuccessRate   float64   `gorm:"column:success_rate" json:"success_rate"`
	FailureRate   float64   `gorm:"column:failure_rate" json:"failure_rate"`
	P50DurationMs float64   `gorm:"column:p50_duration_ms" json:"p50_duration_ms"`
	P95DurationMs float64   `gorm:"column:p95_duration_ms" json:"p95_duration_ms"`
	P99DurationMs float64   `gorm:"column:p99_duration_ms" json:"p99_duration_ms"`
	CreatedAt     int64     `gorm:"column:created_at" json:"created_at"`
}

func (EntityMetric) TableName() string { return "entity_metrics" }

type EntityTransition struct {
	ID                   int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	FromEntityID         int64     `gorm:"not null;index:idx_entity_transitions_from_entity_id" json:"from_entity_id"`
	ToEntityID           int64     `gorm:"not null;index:idx_entity_transitions_to_entity_id" json:"to_entity_id"`
	ConditionDescription string    `gorm:"type:text" json:"condition_description"`
	ConditionExpression  string    `gorm:"type:text" json:"condition_expression"`
	TransitionType       string    `gorm:"type:varchar(128);index:idx_entity_transitions_transition_type" json:"transition_type"`
	CreatedAt            time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	FromEntity           Entity    `gorm:"foreignKey:FromEntityID" json:"from_entity,omitempty"`
	ToEntity             Entity    `gorm:"foreignKey:ToEntityID" json:"to_entity,omitempty"`
}

func (EntityTransition) TableName() string { return "entity_transitions" }

type Entity struct {
	ID           int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	FeatureID    *int64         `gorm:"default:NULL;index:idx_entities_feature_id" json:"feature_id"`
	Code         string         `gorm:"type:varchar(128);not null;uniqueIndex" json:"code"`
	Name         string         `gorm:"type:varchar(128);not null;index:idx_entities_name" json:"name"`
	Description  string         `gorm:"type:text" json:"description"`
	DisplayOrder int            `gorm:"default:0" json:"display_order"`
	IsStart      bool           `gorm:"default:false;index:idx_entities_is_start" json:"is_start"`
	IsTerminal   bool           `gorm:"default:false;index:idx_entities_is_terminal" json:"is_terminal"`
	CreatedAt    time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"default:NULL" json:"deleted_at"`
	Feature      Feature        `gorm:"foreignKey:FeatureID" json:"feature,omitempty"`
}

func (Entity) TableName() string { return "entities" }

type FeatureInstance struct {
	ID          string `gorm:"column:id" json:"id"`
	FeatureID   int64  `gorm:"column:feature_id" json:"feature_id"`
	Status      string `gorm:"column:status;type:LowCardinality(String)" json:"status"`
	StartedAt   int64  `gorm:"column:started_at" json:"started_at"`
	CompletedAt *int64 `gorm:"column:completed_at" json:"completed_at,omitempty"`
	CreatedAt   int64  `gorm:"column:created_at" json:"created_at"`
}

func (FeatureInstance) TableName() string { return "feature_instances" }

type FeatureMetric struct {
	FeatureID     int64     `gorm:"column:feature_id" json:"feature_id"`
	WindowStart   time.Time `gorm:"column:window_start" json:"window_start"`
	WindowMinutes int32     `gorm:"column:window_minutes" json:"window_minutes"`
	TotalCount    int64     `gorm:"column:total_count" json:"total_count"`
	SuccessRate   float64   `gorm:"column:success_rate" json:"success_rate"`
	FailureRate   float64   `gorm:"column:failure_rate" json:"failure_rate"`
	P50DurationMs float64   `gorm:"column:p50_duration_ms" json:"p50_duration_ms"`
	P95DurationMs float64   `gorm:"column:p95_duration_ms" json:"p95_duration_ms"`
	P99DurationMs float64   `gorm:"column:p99_duration_ms" json:"p99_duration_ms"`
	CreatedAt     int64     `gorm:"column:created_at" json:"created_at"`
}

func (FeatureMetric) TableName() string { return "feature_metrics" }

type FeatureTeamRole struct {
	FeatureID  int64  `gorm:"column:feature_id" json:"feature_id"`
	TeamID     int64  `gorm:"column:team_id" json:"team_id"`
	Role       string `gorm:"column:role;type:LowCardinality(String)" json:"role"`
	AssignedAt int64  `gorm:"column:assigned_at" json:"assigned_at"`
}

// func (FeatureTeamRole) TableName() string { return "feature_team_roles" }

type Feature struct {
	ID          int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Code        string         `gorm:"type:varchar(128);not null;uniqueIndex" json:"code"`
	Name        string         `gorm:"type:varchar(128);not null;index:idx_features_name" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	IsActive    bool           `gorm:"default:true;index:idx_features_is_active" json:"is_active"`
	CreatedAt   time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"default:NULL" json:"deleted_at"`
}

func (Feature) TableName() string { return "features" }

type KpiRelationship struct {
	ID               int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	KpiID            int64      `gorm:"not null;index:idx_kpis_relationships_kpi_id" json:"kpi_id"`
	RelationType     string     `gorm:"type:varchar(128);not null;index:idx_kpis_relationships_relation_type" json:"relation_type"`
	TargetType       string     `gorm:"type:varchar(128);not null;index:idx_kpis_relationships_target" json:"target_type"`
	TargetID         int64      `gorm:"not null;index:idx_kpis_relationships_target" json:"target_id"`
	Weight           float64    `gorm:"default:0" json:"weight"`
	WeightSetBy      string     `gorm:"type:varchar(128)" json:"weight_set_by"`
	WeightReviewedBy string     `gorm:"type:varchar(128)" json:"weight_reviewed_by"`
	WeightReviewedAt *time.Time `gorm:"default:NULL" json:"weight_reviewed_at,omitempty"`
	CreatedAt        time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	Kpi              Kpi        `gorm:"foreignKey:KpiID" json:"kpi,omitempty"`
}

func (KpiRelationship) TableName() string { return "kpis_relationships" }

type Kpi struct {
	ID               int64             `gorm:"primaryKey;autoIncrement" json:"id"`
	Code             string            `gorm:"type:varchar(128);not null;uniqueIndex" json:"code"`
	Name             string            `gorm:"type:varchar(128);not null;index:idx_kpis_name" json:"name"`
	Description      string            `gorm:"type:text" json:"description"`
	MetricType       string            `gorm:"type:varchar(128);not null;index:idx_kpis_metric_type" json:"metric_type"`
	Unit             string            `gorm:"type:varchar(128)" json:"unit"`
	CreatedAt        time.Time         `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        time.Time         `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	KpiRelationships []KpiRelationship `gorm:"foreignKey:KpiID" json:"kpi_relationships,omitempty"`
}

func (Kpi) TableName() string { return "kpis" }

type ServiceDeployment struct {
	ID          int64  `gorm:"column:id" json:"id"`
	ServiceID   int64  `gorm:"column:service_id" json:"service_id"`
	Environment string `gorm:"column:environment;type:LowCardinality(String)" json:"environment"`
	CommitHash  string `gorm:"column:commit_hash" json:"commit_hash"`
	Version     string `gorm:"column:version" json:"version"`
	DeployedAt  int64  `gorm:"column:deployed_at" json:"deployed_at"`
	CreatedAt   int64  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   int64  `gorm:"column:updated_at" json:"updated_at"`
}

func (ServiceDeployment) TableName() string { return "service_deployments" }

type ServiceTeamRole struct {
	ID         int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	ServiceID  int64      `gorm:"not null;index:idx_service_team_roles_service_id" json:"service_id"`
	TeamID     int64      `gorm:"not null;index:idx_service_team_roles_team_id" json:"team_id"`
	Role       string     `gorm:"type:varchar(128);not null;index:idx_service_team_roles_role" json:"role"`
	AssignedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"assigned_at"`
	RevokedAt  *time.Time `gorm:"default:NULL" json:"revoked_at,omitempty"`
	Service    Service    `gorm:"foreignKey:ServiceID" json:"service,omitempty"`
	Team       Team       `gorm:"foreignKey:TeamID" json:"team,omitempty"`
}

func (ServiceTeamRole) TableName() string { return "service_team_roles" }

type Service struct {
	ID               int64             `gorm:"primaryKey;autoIncrement" json:"id"`
	Code             string            `gorm:"type:varchar(128);not null;uniqueIndex" json:"code"`
	Name             string            `gorm:"type:varchar(128);not null;index:idx_services_name" json:"name"`
	Description      string            `gorm:"type:text" json:"description"`
	RepositoryURL    string            `gorm:"type:varchar(512)" json:"repository_url"`
	CriticalityLevel string            `gorm:"type:varchar(128);index:idx_services_criticality_level" json:"criticality_level"`
	CreatedAt        time.Time         `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        time.Time         `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt        gorm.DeletedAt    `gorm:"default:NULL" json:"deleted_at"`
	ServiceTeamRoles []ServiceTeamRole `gorm:"foreignKey:ServiceID" json:"service_team_roles,omitempty"`
}

func (Service) TableName() string { return "services" }

type Team struct {
	ID               int64             `gorm:"primaryKey;autoIncrement" json:"id"`
	Name             string            `gorm:"type:varchar(128);not null;index:idx_teams_name" json:"name"`
	SlackChannel     string            `gorm:"type:varchar(128)" json:"slack_channel"`
	PagerdutyKey     string            `gorm:"type:varchar(128)" json:"pagerduty_key"`
	OncallEmail      string            `gorm:"type:varchar(255)" json:"oncall_email"`
	CreatedAt        time.Time         `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        time.Time         `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt        gorm.DeletedAt    `gorm:"default:NULL" json:"deleted_at"`
	ServiceTeamRoles []ServiceTeamRole `gorm:"foreignKey:TeamID" json:"service_team_roles,omitempty"`
}

func (Team) TableName() string { return "teams" }

// --- Seeding Constants ---
const (
	numTeams              = 20
	numFeatures           = 20
	numServices           = 20
	numApis               = 20
	numKpis               = 20
	numEntities           = 20
	numServiceDeployments = 20
	numApiDeployments     = 20
	numApiExecutions      = 20
	numApiDbInteractions  = 20
	numKpiRelationships   = 20
	numServiceTeamRoles   = 20
	// numFeatureTeamRoles   = 20
	numEntityApis        = 20
	numEntityTransitions = 20
	numFeatureInstances  = 20
	numEntityInstances   = 20
	numApiMetrics        = 20
	numEntityMetrics     = 20
	numFeatureMetrics    = 20
)

var (
	serviceTiers      = []string{"P0", "P1", "P2"}
	deploymentEnvs    = []string{"dev", "staging", "prod"}
	featureStatuses   = []string{"PENDING", "ACTIVE", "COMPLETED", "FAILED"}
	entityStatuses    = []string{"PENDING", "RUNNING", "COMPLETED", "FAILED"}
	serviceRoleTypes  = []string{"OWNER", "CONTRIBUTOR", "ONCALL"}
	kpiRelationTypes  = []string{"impact", "dependency", "ownership"}
	transitionTypes   = []string{"AUTOMATIC", "MANUAL", "EVENT_DRIVEN"}
	databaseTypes     = []string{"PostgreSQL", "ClickHouse", "Redis"}
	dbOperationTypes  = []string{"SELECT", "INSERT", "UPDATE"}
	apiEntitySeedData []EntityApi
	featureSeedData   = []featureSeed{
		{name: "customer_onboarding", description: "Starts and tracks retail customer onboarding from identity capture through account activation."},
		{name: "agent_signin", description: "Handles agent authentication, session creation, and sign-in risk checks."},
		{name: "forget_mpin", description: "Resets a user's MPIN after identity verification and OTP confirmation."},
		{name: "change_mpin", description: "Lets an authenticated user rotate their MPIN with policy validation."},
		{name: "internal_transfer", description: "Processes wallet-to-wallet transfers between accounts within the platform."},
		{name: "order_payment", description: "Authorizes and records payment collection for an existing order."},
		{name: "order_placement", description: "Creates a new order and reserves downstream inventory and payment state."},
		{name: "email_forget_password", description: "Initiates email-based password recovery and recovery token validation."},
		{name: "phone_otp_signup", description: "Registers a new user account using phone verification by OTP."},
		{name: "email_verify", description: "Verifies a user's email address and updates the account verification state."},
		{name: "agent_forget_password", description: "Starts password reset for support agents with additional access controls."},
		{name: "agent_change_password_by_otp", description: "Allows agents to change password after OTP-based verification."},
		{name: "agent_onboarding", description: "Provisions a new agent profile, role assignments, and initial credentials."},
		{name: "agent_change_password_by_magic_link", description: "Changes an agent password using a signed magic-link recovery flow."},
		{name: "agent_unlock_user", description: "Unlocks a blocked user account through authorized agent intervention."},
		{name: "merchant_onboarding", description: "Onboards a merchant profile, business details, and settlement configuration."},
		{name: "merchant_forget_mpin", description: "Resets merchant MPIN after ownership validation and recovery approval."},
		{name: "unblock_user_by_agent", description: "Lets an agent remove temporary restrictions from a user account."},
		{name: "reset_user_pin_by_agent", description: "Resets a user transaction PIN through the assisted agent workflow."},
		{name: "external_transfer", description: "Submits and monitors transfers from the platform to an external bank endpoint."},
	}
	apiSeedData = []apiSeed{
		{endpoint: "/workflows/v1/initiate", method: "POST", description: "Initiates workflow execution for onboarding, recovery, transfer, and order journeys."},
		{endpoint: "/notification/v1/send", method: "POST", description: "Sends transactional notifications for OTP, recovery, and workflow status updates."},
		{endpoint: "/kyc/v1/submit", method: "POST", description: "Submits identity documents and profile payloads for KYC verification."},
		{endpoint: "/notification/v1/template", method: "GET", description: "Fetches notification template metadata used before sending user or agent alerts."},
		{endpoint: "/notification/v1/client", method: "POST", description: "Registers or updates a notification client used by workflow channels."},
		{endpoint: "/alert/v1/channels", method: "GET", description: "Lists alert delivery channels configured for operational or risk notifications."},
		{endpoint: "/users/v1/roles/create", method: "POST", description: "Creates a new user or agent role for workflow authorization rules."},
		{endpoint: "/users/v1/roles/update", method: "PUT", description: "Updates an existing user or agent role definition."},
		{endpoint: "/users/v1/attachment", method: "POST", description: "Uploads user attachments such as KYC or onboarding documents."},
		{endpoint: "/otp/private/v1/client", method: "POST", description: "Creates an OTP client configuration for internal verification flows."},
		{endpoint: "/otp/private/v1/verify", method: "POST", description: "Verifies an OTP for signup, password reset, or MPIN recovery."},
		{endpoint: "/users/v1/profile", method: "POST", description: "Creates or updates personal profile details collected during onboarding."},
		{endpoint: "/users/v1/password/reset", method: "POST", description: "Resets a user password after recovery token validation."},
		{endpoint: "/users/v1/mpin/reset", method: "POST", description: "Resets a user or merchant MPIN after OTP or agent approval."},
		{endpoint: "/users/v1/pin/reset", method: "POST", description: "Resets a transaction PIN through an assisted support workflow."},
		{endpoint: "/agent/v1/auth/signin", method: "POST", description: "Authenticates an agent and starts an agent console session."},
		{endpoint: "/agent/v1/users/unlock", method: "POST", description: "Unlocks a blocked user account through agent action."},
		{endpoint: "/merchant/v1/onboarding", method: "POST", description: "Submits merchant onboarding details and starts merchant review."},
		{endpoint: "/transfer/v1/internal", method: "POST", description: "Executes an internal transfer between platform accounts."},
		{endpoint: "/transfer/v1/external", method: "POST", description: "Executes an external transfer to a bank or partner endpoint."},
	}
	entitySeedData = []entitySeed{
		{name: "onboarding_start", description: "Entry state for customer onboarding before any verification begins."},
		{name: "verify_otp", description: "Validates the OTP submitted during onboarding or recovery."},
		{name: "set_pin", description: "Collects and stores the initial user PIN during onboarding."},
		{name: "verify_pin", description: "Confirms the user-provided PIN before allowing workflow progression."},
		{name: "set_personal_details", description: "Captures customer demographic and profile details."},
		{name: "order_placement_start", description: "Initial state for a new order placement request."},
		{name: "verify_mpin_for_order_placement", description: "Verifies MPIN before allowing the order placement request to continue."},
		{name: "order_placement_submit", description: "Submits the final order placement payload to downstream services."},
		{name: "order_placement_submit_alt", description: "Alternate order placement submit state used by the frontend to keep the user on the MPIN screen."},
		{name: "change_mpin_start", description: "Starts the authenticated change MPIN flow."},
		{name: "forget_mpin_start", description: "Starts the forgot MPIN recovery journey."},
		{name: "verify_otp_and_generate_reset_token", description: "Verifies OTP and generates a secure reset token."},
		{name: "protected_reset_mpin", description: "Performs the protected MPIN reset using a verified token."},
		{name: "agent_signin_start", description: "Initial state for agent sign-in before secondary verification."},
		{name: "verify_otp_get_token", description: "Validates OTP and issues an access token to the agent."},
		{name: "email_forget_password_start", description: "Starts email-based password recovery."},
		{name: "verify_email_otp", description: "Verifies the OTP sent to the user's email address."},
		{name: "reset_password", description: "Resets the password after the recovery checks have passed."},
		{name: "phone_otp_signup_start", description: "Starts signup with phone number based OTP verification."},
		{name: "verify_otp_and_check_user", description: "Verifies OTP and checks whether the user already exists."},
	}
	serviceSeedData = []namedDescriptionSeed{
		{name: "Workflow Service", description: "Coordinates workflow state changes and execution routing across the platform."},
		{name: "Notification Service", description: "Sends OTP, email, and in-app notifications for user-facing workflow events."},
		{name: "KYC Service", description: "Integrates with kyc verification providers and stores KYC submission outcomes."},
		{name: "Users Service", description: "Manages user credentials, recovery flows, and authentication checks."},
		{name: "Payments Service", description: "Executes internal and external transfer instructions with settlement controls."},
		{name: "Order Service", description: "Creates orders, tracks status, and coordinates payment and fulfillment events."},
		{name: "Ontology Service", description: "Analyses whole system, onboarding status along with all insights."},
		{name: "Agent Console API", description: "Supports agent actions such as unlock, PIN reset, and assisted recovery."},
		{name: "Risk Engine", description: "Applies fraud, velocity, and policy rules to high-risk workflow events."},
		{name: "Audit Service", description: "Stores audit logs and compliance trails for sensitive workflow operations."},
		{name: "Auth Service", description: "Maintains user demographic and contact profile information."},
		{name: "Kong Service", description: "API gateway."},
		{name: "Kyb Service", description: "Handles password, MPIN, and PIN lifecycle operations."},
		{name: "Payment Gateway", description: "Connects order and transfer flows to external payment providers."},
		{name: "Ledger Service", description: "Records financial postings generated by payment and transfer workflows."},
		{name: "Compliance Service", description: "Evaluates AML, sanctions, and regulatory controls for workflows."},
		{name: "OTP Service", description: "Generates and validates one-time passwords across signup and recovery flows."},
		{name: "Email Service", description: "Sends verification and recovery emails with delivery tracking."},
		{name: "Customer Registry", description: "Stores master customer records used by onboarding and servicing workflows."},
		{name: "Merchant Registry", description: "Stores merchant master data and onboarding lifecycle details."},
	}
	teamSeedData = []teamSeed{
		{name: "Workflow Platform", slack: "#workflow-platform", oncall: "workflow-platform-oncall@example.com"},
		{name: "Identity Ops", slack: "#identity-ops", oncall: "identity-ops-oncall@example.com"},
		{name: "Transfer Reliability", slack: "#transfer-reliability", oncall: "transfer-reliability-oncall@example.com"},
		{name: "Merchant Success", slack: "#merchant-success", oncall: "merchant-success-oncall@example.com"},
		{name: "Agent Experience", slack: "#agent-experience", oncall: "agent-experience-oncall@example.com"},
		{name: "Risk Controls", slack: "#risk-controls", oncall: "risk-controls-oncall@example.com"},
		{name: "Notification Delivery", slack: "#notification-delivery", oncall: "notification-delivery-oncall@example.com"},
		{name: "Compliance Engineering", slack: "#compliance-engineering", oncall: "compliance-engineering-oncall@example.com"},
		{name: "Order Fulfillment", slack: "#order-fulfillment", oncall: "order-fulfillment-oncall@example.com"},
		{name: "Payments Core", slack: "#payments-core", oncall: "payments-core-oncall@example.com"},
		{name: "Customer Identity", slack: "#customer-identity", oncall: "customer-identity-oncall@example.com"},
		{name: "Merchant Identity", slack: "#merchant-identity", oncall: "merchant-identity-oncall@example.com"},
		{name: "Audit And Governance", slack: "#audit-governance", oncall: "audit-governance-oncall@example.com"},
		{name: "Session Security", slack: "#session-security", oncall: "session-security-oncall@example.com"},
		{name: "Profile Platform", slack: "#profile-platform", oncall: "profile-platform-oncall@example.com"},
		{name: "KYC Operations", slack: "#kyc-operations", oncall: "kyc-operations-oncall@example.com"},
		{name: "OTP Platform", slack: "#otp-platform", oncall: "otp-platform-oncall@example.com"},
		{name: "Email Platform", slack: "#email-platform", oncall: "email-platform-oncall@example.com"},
		{name: "Ledger Reliability", slack: "#ledger-reliability", oncall: "ledger-reliability-oncall@example.com"},
		{name: "Customer Support Tools", slack: "#support-tools", oncall: "support-tools-oncall@example.com"},
	}
	kpiSeedData = []kpiSeed{
		{name: "Workflow Initiation Latency", description: "Measures time to accept and enqueue workflow initiation requests.", metricType: "Latency", unit: "ms"},
		{name: "Notification Delivery Rate", description: "Tracks successful delivery ratio for OTP and transactional notifications.", metricType: "Availability", unit: "percent"},
		{name: "KYC Submission Error Rate", description: "Measures submission failures returned by the KYC processing pipeline.", metricType: "ErrorRate", unit: "percent"},
		{name: "Agent Sign-In Throughput", description: "Counts successful agent authentication requests processed per window.", metricType: "Throughput", unit: "count"},
		{name: "Password Recovery Completion Rate", description: "Tracks completed password recovery journeys versus initiated ones.", metricType: "Availability", unit: "percent"},
		{name: "MPIN Reset Latency", description: "Measures how quickly MPIN reset requests complete after verification.", metricType: "Latency", unit: "ms"},
		{name: "Internal Transfer Success Rate", description: "Tracks successful completion ratio for internal transfer workflows.", metricType: "Availability", unit: "percent"},
		{name: "External Transfer Error Rate", description: "Measures failed external transfer executions and downstream rejections.", metricType: "ErrorRate", unit: "percent"},
		{name: "Order Placement Throughput", description: "Counts order placement requests completed during the metrics window.", metricType: "Throughput", unit: "count"},
		{name: "Order Payment Latency", description: "Measures processing time for order payment authorization and capture.", metricType: "Latency", unit: "ms"},
		{name: "Email Verification Success Rate", description: "Tracks successful email verification completions for pending accounts.", metricType: "Availability", unit: "percent"},
		{name: "Phone Signup Error Rate", description: "Measures signup failures in the phone OTP registration flow.", metricType: "ErrorRate", unit: "percent"},
		{name: "Agent Onboarding Throughput", description: "Counts newly provisioned agent accounts completed within the window.", metricType: "Throughput", unit: "count"},
		{name: "Merchant Onboarding Latency", description: "Measures elapsed time from merchant submission to onboarding completion.", metricType: "Latency", unit: "ms"},
		{name: "Merchant MPIN Recovery Rate", description: "Tracks successful merchant MPIN recovery outcomes.", metricType: "Availability", unit: "percent"},
		{name: "User Unlock Error Rate", description: "Measures failures during agent-assisted user unlock operations.", metricType: "ErrorRate", unit: "percent"},
		{name: "PIN Reset Throughput", description: "Counts completed user PIN reset actions performed by agents.", metricType: "Throughput", unit: "count"},
		{name: "Risk Review Latency", description: "Measures turnaround time for workflows routed through risk review.", metricType: "Latency", unit: "ms"},
		{name: "Audit Logging Availability", description: "Tracks successful persistence of audit records for critical actions.", metricType: "Availability", unit: "percent"},
		{name: "Ledger Posting Error Rate", description: "Measures failures when recording postings for payment and transfer flows.", metricType: "ErrorRate", unit: "percent"},
	}
)

type featureSeed struct {
	name        string
	description string
}

type apiSeed struct {
	endpoint    string
	method      string
	description string
}

type entitySeed struct {
	name        string
	description string
}

type namedDescriptionSeed struct {
	name        string
	description string
}

type teamSeed struct {
	name   string
	slack  string
	oncall string
}

type kpiSeed struct {
	name        string
	description string
	metricType  string
	unit        string
}

// --- Main Function ---

func main() {
	var flagConfig = flag.String("config", "./config/local.yml", "path to the config file")
	flag.Parse()

	cfg, err := config.Load(*flagConfig)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	seedTarget := strings.ToLower(strings.TrimSpace(os.Getenv("SEED_TARGET")))
	if seedTarget == "" {
		seedTarget = "both"
	}
	if seedTarget != "both" && seedTarget != "postgres" && seedTarget != "clickhouse" {
		log.Fatalf("Unsupported SEED_TARGET %q. Allowed values: postgres, clickhouse, both", seedTarget)
	}

	seedPostgres := seedTarget != "clickhouse"
	seedClickHouseRequested := seedTarget != "postgres"
	if seedClickHouseRequested && !cfg.ClickHouse.Enabled {
		log.Fatalf("ClickHouse seed requested via SEED_TARGET=%s but ClickHouse is disabled in config", seedTarget)
	}

	fmt.Printf("Seed target: %s\n", seedTarget)

	// fmt.Printf("Loaded DSN: '%s'\n", cfg.Database.MasterDatabaseDsn)

	pgDB, err := gorm.Open(postgres.Open(cfg.Database.MasterDatabaseDsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to postgres: %v", err)
	}

	var clickhouseDB *gorm.DB
	if seedClickHouseRequested && cfg.ClickHouse.Enabled {
		clickhouseDB, err = gorm.Open(clickhouse.Open(cfg.ClickHouse.DSN), &gorm.Config{})
		if err != nil {
			log.Fatalf("Failed to connect to clickhouse: %v", err)
		}
	}

	fmt.Println("🚀 Starting Full System Seed...")

	var (
		teams    []Team
		features []Feature
		services []Service
		apis     []Api
		kpis     []Kpi
		entities []Entity
	)

	if seedPostgres {
		teams, err = seedTeams(pgDB)
		if err != nil {
			log.Fatalf("Failed to seed teams: %v", err)
		}
		features, err = seedFeatures(pgDB)
		if err != nil {
			log.Fatalf("Failed to seed features: %v", err)
		}
		services, err = seedServices(pgDB)
		if err != nil {
			log.Fatalf("Failed to seed services: %v", err)
		}
		apis, err = seedApis(pgDB)
		if err != nil {
			log.Fatalf("Failed to seed apis: %v", err)
		}
		kpis, err = seedKpis(pgDB)
		if err != nil {
			log.Fatalf("Failed to seed kpis: %v", err)
		}
		entities, err = seedEntities(pgDB, features)
		if err != nil {
			log.Fatalf("Failed to seed entities: %v", err)
		}
	} else {
		fmt.Println("Skipping Postgres master data seeding; reusing existing records for reference data.")
		features, err = loadFeatures(pgDB)
		if err != nil {
			log.Fatalf("Failed to load features from postgres: %v", err)
		}
		services, err = loadServices(pgDB)
		if err != nil {
			log.Fatalf("Failed to load services from postgres: %v", err)
		}
		apis, err = loadApis(pgDB)
		if err != nil {
			log.Fatalf("Failed to load apis from postgres: %v", err)
		}
		entities, err = loadEntities(pgDB)
		if err != nil {
			log.Fatalf("Failed to load entities from postgres: %v", err)
		}
	}

	if seedPostgres {
		if err := seedKpiRelationships(pgDB, kpis, services, features, entities, apis); err != nil {
			log.Fatalf("Failed to seed kpi relationships: %v", err)
		}
		if err := seedServiceTeamRoles(pgDB, services, teams); err != nil {
			log.Fatalf("Failed to seed service team roles: %v", err)
		}
	} else {
		fmt.Println("Skipping Postgres-only relationships and team role seeding.")
	}

	deploymentsDB := pgDB
	if clickhouseDB != nil {
		deploymentsDB = clickhouseDB
	}
	deployments, err := seedServiceDeployments(deploymentsDB, services)
	if err != nil {
		log.Fatalf("Failed to seed service deployments: %v", err)
	}

	if seedPostgres {
		apiEntitySeedData, err = seedEntityApis(pgDB, entities, apis)
		if err != nil {
			log.Fatalf("Failed to seed entity apis: %v", err)
		}
		if err := seedEntityTransitions(pgDB, entities); err != nil {
			log.Fatalf("Failed to seed entity transitions: %v", err)
		}
	} else {
		fmt.Println("Skipping Postgres-only entity APIs and transitions.")
	}

	apiDeploymentsDB := pgDB
	if clickhouseDB != nil {
		apiDeploymentsDB = clickhouseDB
	}
	apiDeployments, err := seedApiDeployments(apiDeploymentsDB, apis, deployments)
	if err != nil {
		log.Fatalf("Failed to seed api deployments: %v", err)
	}

	featureInstancesDB := pgDB
	if clickhouseDB != nil {
		featureInstancesDB = clickhouseDB
	}
	featureInstances, err := seedFeatureInstances(featureInstancesDB, features)
	if err != nil {
		log.Fatalf("Failed to seed feature instances: %v", err)
	}

	entityInstancesDB := pgDB
	if clickhouseDB != nil {
		entityInstancesDB = clickhouseDB
	}
	entityInstances, err := seedEntityInstances(entityInstancesDB, featureInstances, entities)
	if err != nil {
		log.Fatalf("Failed to seed entity instances: %v", err)
	}

	apiExecutionsDB := pgDB
	if clickhouseDB != nil {
		apiExecutionsDB = clickhouseDB
	}
	apiExecutions, err := seedApiExecutions(apiExecutionsDB, entityInstances, apiDeployments)
	if err != nil {
		log.Fatalf("Failed to seed api executions: %v", err)
	}

	apiDbInteractionsDB := pgDB
	if clickhouseDB != nil {
		apiDbInteractionsDB = clickhouseDB
	}
	if err := seedApiDbInteractions(apiDbInteractionsDB, apiExecutions); err != nil {
		log.Fatalf("Failed to seed api db interactions: %v", err)
	}

	apiMetricsDB := pgDB
	if clickhouseDB != nil {
		apiMetricsDB = clickhouseDB
	}
	if err := seedApiMetrics(apiMetricsDB, apiDeployments); err != nil {
		log.Fatalf("Failed to seed api metrics: %v", err)
	}

	entityMetricsDB := pgDB
	if clickhouseDB != nil {
		entityMetricsDB = clickhouseDB
	}
	if err := seedEntityMetrics(entityMetricsDB, entities); err != nil {
		log.Fatalf("Failed to seed entity metrics: %v", err)
	}

	featureMetricsDB := pgDB
	if clickhouseDB != nil {
		featureMetricsDB = clickhouseDB
	}
	if err := seedFeatureMetrics(featureMetricsDB, features); err != nil {
		log.Fatalf("Failed to seed feature metrics: %v", err)
	}

	fmt.Println("✅ Seeding Complete! All tables populated.")
}

// --- HELPER FUNCTIONS ---

func seedTeams(db *gorm.DB) ([]Team, error) {
	fmt.Println("Seeding teams...")
	var list []Team
	for i := 1; i <= numTeams; i++ {
		team := teamSeedData[i-1]
		list = append(list, Team{
			Name:         team.name,
			SlackChannel: team.slack,
			PagerdutyKey: fmt.Sprintf("pd_key_%03d", i),
			OncallEmail:  team.oncall,
		})
	}
	if err := db.Create(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func seedFeatures(db *gorm.DB) ([]Feature, error) {
	fmt.Println("Seeding features...")
	var list []Feature
	for i := 1; i <= numFeatures; i++ {
		workflow := featureSeedData[i-1]
		list = append(list, Feature{
			Code:        strings.ToLower(strings.ReplaceAll(workflow.name, " ", "_")),
			Name:        workflow.name,
			Description: workflow.description,
			IsActive:    i%9 != 0,
		})
	}
	if err := db.Create(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func seedServices(db *gorm.DB) ([]Service, error) {
	fmt.Println("Seeding services...")
	var list []Service
	for i := 1; i <= numServices; i++ {
		service := serviceSeedData[i-1]
		list = append(list, Service{
			Code:             strings.ToLower(strings.ReplaceAll(service.name, " ", "_")),
			Name:             service.name,
			Description:      service.description,
			RepositoryURL:    fmt.Sprintf("https://git.example.com/services/%s.git", strings.ToLower(strings.ReplaceAll(service.name, " ", "-"))),
			CriticalityLevel: serviceTiers[(i-1)%len(serviceTiers)],
		})
	}
	if err := db.Create(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func seedApis(db *gorm.DB) ([]Api, error) {
	fmt.Println("Seeding apis...")
	var list []Api
	for i := 1; i <= numApis; i++ {
		api := apiSeedData[(i-1)%len(apiSeedData)]
		list = append(list, Api{
			Endpoint:    api.endpoint,
			HttpMethod:  api.method,
			Protocol:    "HTTP",
			IsInternal:  strings.HasPrefix(api.endpoint, "/otp/private/"),
			Description: fmt.Sprintf("%s Workflow mapping: %s.", api.description, featureSeedData[(i-1)%len(featureSeedData)].name),
		})
	}
	if err := db.Create(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func seedEntities(db *gorm.DB, features []Feature) ([]Entity, error) {
	fmt.Println("Seeding entities...")
	var list []Entity
	for i := 1; i <= numEntities; i++ {
		f := features[(i-1)%len(features)]
		entitySeed := entitySeedData[i-1]
		list = append(list, Entity{
			FeatureID:    &f.ID,
			Code:         strings.ToLower(strings.ReplaceAll(entitySeed.name, " ", "_")),
			Name:         entitySeed.name,
			Description:  fmt.Sprintf("%s Workflow: %s.", entitySeed.description, f.Description),
			DisplayOrder: i,
			IsStart:      strings.HasSuffix(entitySeed.name, "_start"),
			IsTerminal:   strings.Contains(entitySeed.name, "submit") || strings.HasPrefix(entitySeed.name, "reset_") || strings.HasPrefix(entitySeed.name, "protected_reset"),
		})
	}
	if err := db.Create(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func seedServiceDeployments(db *gorm.DB, svcs []Service) ([]ServiceDeployment, error) {
	fmt.Println("Seeding service deployments...")
	var list []ServiceDeployment
	for i := 1; i <= numServiceDeployments; i++ {
		s := svcs[(i-1)%len(svcs)]
		deployedAt := time.Now().Add(-time.Duration(i) * 6 * time.Hour).Unix()
		list = append(list, ServiceDeployment{
			ID:          int64(i),
			ServiceID:   s.ID,
			Environment: deploymentEnvs[(i-1)%len(deploymentEnvs)],
			CommitHash:  fmt.Sprintf("%040x", i*7919),
			Version:     fmt.Sprintf("v1.%d.%d", (i-1)/10, i%10),
			DeployedAt:  deployedAt,
			CreatedAt:   deployedAt,
			UpdatedAt:   deployedAt,
		})
	}
	if err := db.Create(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func seedApiDeployments(db *gorm.DB, apis []Api, deploys []ServiceDeployment) ([]ApisDeployment, error) {
	fmt.Println("Seeding api deployments...")
	var list []ApisDeployment
	for i := 1; i <= numApiDeployments; i++ {
		api := apis[(i-1)%len(apis)]
		deployment := deploys[(i-1)%len(deploys)]
		list = append(list, ApisDeployment{
			ID:                  uuid.NewString(),
			ApiID:               api.ID,
			ServiceDeploymentID: deployment.ID,
		})
	}
	if err := db.Create(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func seedApiExecutions(db *gorm.DB, ei []EntityInstance, ad []ApisDeployment) ([]ApiExecution, error) {
	fmt.Println("Seeding api executions...")
	var list []ApiExecution
	deploymentsByAPI := make(map[int64][]ApisDeployment)
	for _, deployment := range ad {
		deploymentsByAPI[deployment.ApiID] = append(deploymentsByAPI[deployment.ApiID], deployment)
	}

	entityAPIMap := make(map[int64][]int64)
	for _, relation := range apiEntitySeedData {
		entityAPIMap[relation.EntityID] = append(entityAPIMap[relation.EntityID], relation.ApiID)
	}

	for i := 1; i <= numApiExecutions; i++ {
		entityInstance := ei[(i-1)%len(ei)]
		apiDeployment := ad[(i-1)%len(ad)]
		if apiIDs := entityAPIMap[entityInstance.EntityID]; len(apiIDs) > 0 {
			apiID := apiIDs[(i-1)%len(apiIDs)]
			if candidateDeployments := deploymentsByAPI[apiID]; len(candidateDeployments) > 0 {
				apiDeployment = candidateDeployments[(i-1)%len(candidateDeployments)]
			}
		}

		requestAt := time.Now().Add(-time.Duration(i) * 15 * time.Minute).Unix()
		durationMs := int32(40 + (i % 9 * 25))
		status := int32(200)
		if i%12 == 0 {
			status = 500
		} else if i%5 == 0 {
			status = 202
		}

		sourceIP := fmt.Sprintf("10.0.%d.%d", (i%8)+1, (i%200)+10)
		requestID := uuid.NewString()
		correlationID := uuid.NewString()
		traceID := uuid.NewString()
		spanID := fmt.Sprintf("%016x", i*12345)
		host := fmt.Sprintf("api-node-%02d.internal", (i%6)+1)
		lat := 12.90 + float64(i%10)/100
		lon := 77.50 + float64(i%10)/100
		deviceID := fmt.Sprintf("device-%03d", i)
		platformOS := []string{"ios", "android", "web"}[(i-1)%3]
		appVersion := fmt.Sprintf("2.%d.%d", (i-1)/10, i%10)
		userAgent := fmt.Sprintf("agent-grid-seeder/%s", appVersion)
		riskSessionID := fmt.Sprintf("risk-%03d", i)
		deviceModel := []string{"iPhone", "Pixel", "Chrome"}[(i-1)%3]
		authScheme := []string{"Bearer", "ApiKey"}[(i-1)%2]
		tokenType := []string{"access_token", "service_token"}[(i-1)%2]
		tokenScopes := "read:workflow write:workflow"
		var appErrorCode *string
		var errorMessage *string
		if status >= 400 {
			code := "INTERNAL_ERROR"
			message := "Seeded upstream timeout"
			appErrorCode = &code
			errorMessage = &message
		}

		list = append(list, ApiExecution{
			ID:                uuid.NewString(),
			EntityInstanceID:  entityInstance.ID,
			ApiDeploymentsID:  apiDeployment.ID,
			OauthClientID:     stringPtr(fmt.Sprintf("oauth-client-%02d", (i%10)+1)),
			UserID:            stringPtr(fmt.Sprintf("user-%03d", i)),
			UserUID:           stringPtr(fmt.Sprintf("uid-%03d", i)),
			OrganizationID:    stringPtr(fmt.Sprintf("org-%02d", (i%7)+1)),
			TokenID:           stringPtr(fmt.Sprintf("token-%03d", i)),
			RequestTimestamp:  requestAt,
			RequestDurationMs: durationMs,
			SourceIP:          &sourceIP,
			RequestID:         &requestID,
			CorrelationID:     &correlationID,
			TraceID:           &traceID,
			SpanID:            &spanID,
			Host:              &host,
			Latitude:          &lat,
			Longitude:         &lon,
			DeviceID:          &deviceID,
			PlatformOS:        &platformOS,
			AppVersion:        &appVersion,
			UserAgent:         &userAgent,
			RiskSessionID:     &riskSessionID,
			DeviceModel:       &deviceModel,
			AuthScheme:        &authScheme,
			TokenType:         &tokenType,
			TokenScopes:       &tokenScopes,
			HttpStatus:        status,
			AppErrorCode:      appErrorCode,
			ErrorMessage:      errorMessage,
			ResponseSizeBytes: int32(512 + i*17),
			CreatedAt:         requestAt,
		})
	}
	if err := db.Create(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func seedApiDbInteractions(db *gorm.DB, execs []ApiExecution) error {
	fmt.Println("Seeding api db interactions...")
	var list []ApiDbInteraction
	queryCatalog := []struct {
		tableName string
		query     string
	}{
		{tableName: "features", query: "SELECT * FROM features WHERE code = ?"},
		{tableName: "entities", query: "SELECT * FROM entities WHERE feature_id = ?"},
		{tableName: "apis", query: "SELECT * FROM apis WHERE endpoint = ?"},
		{tableName: "entity_transitions", query: "SELECT * FROM entity_transitions WHERE from_entity_id = ?"},
		{tableName: "entity_apis", query: "INSERT INTO entity_apis(entity_id, api_id) VALUES (?, ?)"},
	}
	for i := 1; i <= numApiDbInteractions; i++ {
		exec := execs[(i-1)%len(execs)]
		dbType := databaseTypes[(i-1)%len(databaseTypes)]
		queryMeta := queryCatalog[(i-1)%len(queryCatalog)]
		operation := dbOperationTypes[(i-1)%len(dbOperationTypes)]
		query := queryMeta.query
		if operation == "UPDATE" && queryMeta.tableName == "entities" {
			query = "UPDATE entities SET description = ? WHERE id = ?"
		}
		params := fmt.Sprintf(`["%s"]`, valueOrDefault(exec.CorrelationID, "unknown-correlation"))
		var errorMessage *string
		if i%11 == 0 {
			msg := "Seeded deadlock detected"
			errorMessage = &msg
		}
		list = append(list, ApiDbInteraction{
			ID:              uuid.NewString(),
			ApiExecutionID:  exec.ID,
			DatabaseType:    dbType,
			DbTableName:     queryMeta.tableName,
			OperationType:   operation,
			SqlQuery:        query,
			QueryParameters: &params,
			RowsAffected:    int32((i % 5) + 1),
			DurationMs:      int32(5 + (i % 7 * 8)),
			ErrorMessage:    errorMessage,
			CreatedAt:       exec.CreatedAt + int64(i%3),
		})
	}
	return db.Create(&list).Error
}

func seedKpis(db *gorm.DB) ([]Kpi, error) {
	fmt.Println("Seeding kpis...")
	var list []Kpi
	for i := 1; i <= numKpis; i++ {
		kpi := kpiSeedData[i-1]
		list = append(list, Kpi{
			Code:        strings.ToLower(strings.ReplaceAll(kpi.name, " ", "_")),
			Name:        kpi.name,
			Description: kpi.description,
			MetricType:  kpi.metricType,
			Unit:        kpi.unit,
		})
	}
	if err := db.Create(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func seedKpiRelationships(db *gorm.DB, kpis []Kpi, services []Service, features []Feature, entities []Entity, apis []Api) error {
	fmt.Println("Seeding kpi relationships...")
	var list []KpiRelationship
	targetGroups := []struct {
		targetType string
		targetIDs  []int64
	}{
		{targetType: "service", targetIDs: collectServiceIDs(services)},
		{targetType: "feature", targetIDs: collectFeatureIDs(features)},
		{targetType: "entity", targetIDs: collectEntityIDs(entities)},
		{targetType: "api", targetIDs: collectAPIIDs(apis)},
	}
	for i := 1; i <= numKpiRelationships; i++ {
		group := targetGroups[(i-1)%len(targetGroups)]
		reviewedAt := time.Now().Add(-time.Duration(i) * time.Hour)
		list = append(list, KpiRelationship{
			KpiID:            kpis[(i-1)%len(kpis)].ID,
			RelationType:     kpiRelationTypes[(i-1)%len(kpiRelationTypes)],
			TargetType:       group.targetType,
			TargetID:         group.targetIDs[(i-1)%len(group.targetIDs)],
			Weight:           0.2 + float64((i%7)+1)/10,
			WeightSetBy:      fmt.Sprintf("owner-%02d", (i%10)+1),
			WeightReviewedBy: fmt.Sprintf("reviewer-%02d", (i%6)+1),
			WeightReviewedAt: &reviewedAt,
		})
	}
	return db.Create(&list).Error
}

func seedServiceTeamRoles(db *gorm.DB, svcs []Service, teams []Team) error {
	fmt.Println("Seeding service team roles...")
	var list []ServiceTeamRole
	for i := 1; i <= numServiceTeamRoles; i++ {
		assignedAt := time.Now().Add(-time.Duration(i) * 24 * time.Hour)
		list = append(list, ServiceTeamRole{
			ServiceID:  svcs[(i-1)%len(svcs)].ID,
			TeamID:     teams[(i+2)%len(teams)].ID,
			Role:       serviceRoleTypes[(i-1)%len(serviceRoleTypes)],
			AssignedAt: assignedAt,
		})
	}
	return db.Create(&list).Error
}

// func seedFeatureTeamRoles(db *gorm.DB, feats []Feature, teams []Team) error {
// 	fmt.Println("Seeding feature team roles...")
// 	var list []FeatureTeamRole
// 	roles := []string{"OWNER", "CONTRIBUTOR", "REVIEWER"}
// 	for i := 1; i <= numFeatureTeamRoles; i++ {
// 		list = append(list, FeatureTeamRole{
// 			FeatureID:  feats[rand.Intn(len(feats))].ID,
// 			TeamID:     teams[rand.Intn(len(teams))].ID,
// 			Role:       roles[rand.Intn(len(roles))],
// 			AssignedAt: time.Now().Unix(),
// 		})
// 	}
// 	return db.Create(&list).Error
// }

func seedEntityApis(db *gorm.DB, ents []Entity, apis []Api) ([]EntityApi, error) {
	fmt.Println("Seeding entity apis...")
	var list []EntityApi
	for i := 1; i <= numEntityApis; i++ {
		list = append(list, EntityApi{
			EntityID: ents[(i-1)%len(ents)].ID,
			ApiID:    apis[(i-1)%len(apis)].ID,
		})
	}
	if err := db.Create(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func loadFeatures(db *gorm.DB) ([]Feature, error) {
	var list []Feature
	if err := db.Find(&list).Error; err != nil {
		return nil, fmt.Errorf("failed to load features: %w", err)
	}
	if len(list) == 0 {
		return nil, fmt.Errorf("no features found to seed; run postgres target first")
	}
	return list, nil
}

func loadServices(db *gorm.DB) ([]Service, error) {
	var list []Service
	if err := db.Find(&list).Error; err != nil {
		return nil, fmt.Errorf("failed to load services: %w", err)
	}
	if len(list) == 0 {
		return nil, fmt.Errorf("no services found to seed; run postgres target first")
	}
	return list, nil
}

func loadApis(db *gorm.DB) ([]Api, error) {
	var list []Api
	if err := db.Find(&list).Error; err != nil {
		return nil, fmt.Errorf("failed to load apis: %w", err)
	}
	if len(list) == 0 {
		return nil, fmt.Errorf("no apis found to seed; run postgres target first")
	}
	return list, nil
}

func loadEntities(db *gorm.DB) ([]Entity, error) {
	var list []Entity
	if err := db.Find(&list).Error; err != nil {
		return nil, fmt.Errorf("failed to load entities: %w", err)
	}
	if len(list) == 0 {
		return nil, fmt.Errorf("no entities found to seed; run postgres target first")
	}
	return list, nil
}

func seedEntityTransitions(db *gorm.DB, ents []Entity) error {
	fmt.Println("Seeding entity transitions...")
	var list []EntityTransition
	entitiesByFeature := make(map[int64][]Entity)
	for _, entity := range ents {
		if entity.FeatureID == nil {
			continue
		}
		entitiesByFeature[*entity.FeatureID] = append(entitiesByFeature[*entity.FeatureID], entity)
	}

	count := 0
	for _, featureEntities := range entitiesByFeature {
		for i := 0; i < len(featureEntities)-1 && count < numEntityTransitions; i++ {
			list = append(list, EntityTransition{
				FromEntityID:         featureEntities[i].ID,
				ToEntityID:           featureEntities[i+1].ID,
				ConditionDescription: fmt.Sprintf("Advance from stage %d to %d when validation passes.", i+1, i+2),
				ConditionExpression:  fmt.Sprintf("entity.display_order == %d && checks_passed == true", i+1),
				TransitionType:       transitionTypes[count%len(transitionTypes)],
			})
			count++
		}
	}

	for i := 1; i <= numEntityTransitions; i++ {
		if len(list) >= numEntityTransitions {
			break
		}
		from := ents[(i-1)%len(ents)]
		to := ents[i%len(ents)]
		if from.ID == to.ID {
			continue
		}
		list = append(list, EntityTransition{
			FromEntityID:         from.ID,
			ToEntityID:           to.ID,
			ConditionDescription: "Fallback transition for seeded non-linear workflow.",
			ConditionExpression:  "retry_count < 3",
			TransitionType:       transitionTypes[(i-1)%len(transitionTypes)],
		})
	}
	return db.Create(&list).Error
}

func seedFeatureInstances(db *gorm.DB, feats []Feature) ([]FeatureInstance, error) {
	fmt.Println("Seeding feature instances...")
	var list []FeatureInstance
	for i := 1; i <= numFeatureInstances; i++ {
		startedAt := time.Now().Add(-time.Duration(i) * time.Hour).Unix()
		status := featureStatuses[(i-1)%len(featureStatuses)]
		var completedAt *int64
		if status == "COMPLETED" || status == "FAILED" {
			ts := startedAt + int64((i%5+1)*600)
			completedAt = &ts
		}
		list = append(list, FeatureInstance{
			ID:          uuid.NewString(),
			FeatureID:   feats[(i-1)%len(feats)].ID,
			Status:      status,
			StartedAt:   startedAt,
			CompletedAt: completedAt,
			CreatedAt:   startedAt,
		})
	}
	if err := db.Create(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func seedEntityInstances(db *gorm.DB, fi []FeatureInstance, ents []Entity) ([]EntityInstance, error) {
	fmt.Println("Seeding entity instances...")
	var list []EntityInstance
	entitiesByFeature := make(map[int64][]Entity)
	for _, entity := range ents {
		if entity.FeatureID == nil {
			continue
		}
		entitiesByFeature[*entity.FeatureID] = append(entitiesByFeature[*entity.FeatureID], entity)
	}
	for i := 1; i <= numEntityInstances; i++ {
		featureInstance := fi[(i-1)%len(fi)]
		candidates := entitiesByFeature[featureInstance.FeatureID]
		if len(candidates) == 0 {
			candidates = ents
		}
		entity := candidates[(i-1)%len(candidates)]
		status := entityStatuses[(i-1)%len(entityStatuses)]
		startedAt := featureInstance.StartedAt + int64(i*60)
		var completedAt *int64
		if status == "COMPLETED" || status == "FAILED" {
			ts := startedAt + int64((i%6+1)*120)
			completedAt = &ts
		}
		list = append(list, EntityInstance{
			ID:                uuid.NewString(),
			FeatureInstanceID: featureInstance.ID,
			EntityID:          entity.ID,
			Status:            status,
			StartedAt:         startedAt,
			CompletedAt:       completedAt,
			CreatedAt:         startedAt,
		})
	}
	if err := db.Create(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func seedApiMetrics(db *gorm.DB, ad []ApisDeployment) error {
	fmt.Println("Seeding api metrics...")
	var list []ApiMetric
	for i := 1; i <= numApiMetrics; i++ {
		windowStart := time.Now().Add(-time.Duration(i) * time.Hour).Truncate(time.Hour)
		successRate := 0.90 + float64(i%8)/100
		errorRate := 1 - successRate
		list = append(list, ApiMetric{
			ApisDeploymentID: ad[(i-1)%len(ad)].ID,
			WindowStart:      windowStart,
			WindowMinutes:    60,
			TotalCalls:       int64(900 + i*25),
			SuccessRate:      successRate,
			ErrorRate:        errorRate,
			P50LatencyMs:     25 + float64(i%10),
			P95LatencyMs:     80 + float64(i%15),
			P99LatencyMs:     120 + float64(i%20),
			CreatedAt:        windowStart.Unix(),
			UpdatedAt:        windowStart.Add(30 * time.Minute).Unix(),
		})
	}
	return db.Create(&list).Error
}

func seedEntityMetrics(db *gorm.DB, ents []Entity) error {
	fmt.Println("Seeding entity metrics...")
	var list []EntityMetric
	for i := 1; i <= numEntityMetrics; i++ {
		windowStart := time.Now().Add(-time.Duration(i) * 2 * time.Hour).Truncate(time.Hour)
		successRate := 0.85 + float64(i%10)/100
		list = append(list, EntityMetric{
			EntityID:      ents[(i-1)%len(ents)].ID,
			WindowStart:   windowStart,
			WindowMinutes: 120,
			TotalCount:    int64(300 + i*12),
			SuccessRate:   successRate,
			FailureRate:   1 - successRate,
			P50DurationMs: 120 + float64(i%12)*5,
			P95DurationMs: 260 + float64(i%15)*8,
			P99DurationMs: 420 + float64(i%18)*10,
			CreatedAt:     windowStart.Unix(),
		})
	}
	return db.Create(&list).Error
}

func seedFeatureMetrics(db *gorm.DB, feats []Feature) error {
	fmt.Println("Seeding feature metrics...")
	var list []FeatureMetric
	for i := 1; i <= numFeatureMetrics; i++ {
		windowStart := time.Now().Add(-time.Duration(i) * 4 * time.Hour).Truncate(time.Hour)
		successRate := 0.88 + float64(i%9)/100
		list = append(list, FeatureMetric{
			FeatureID:     feats[(i-1)%len(feats)].ID,
			WindowStart:   windowStart,
			WindowMinutes: 240,
			TotalCount:    int64(150 + i*9),
			SuccessRate:   successRate,
			FailureRate:   1 - successRate,
			P50DurationMs: 300 + float64(i%10)*11,
			P95DurationMs: 620 + float64(i%10)*19,
			P99DurationMs: 900 + float64(i%12)*24,
			CreatedAt:     windowStart.Unix(),
		})
	}
	return db.Create(&list).Error
}

func stringPtr(v string) *string {
	return &v
}

func valueOrDefault(v *string, fallback string) string {
	if v == nil || *v == "" {
		return fallback
	}
	return *v
}

func collectServiceIDs(services []Service) []int64 {
	ids := make([]int64, 0, len(services))
	for _, service := range services {
		ids = append(ids, service.ID)
	}
	return ids
}

func collectFeatureIDs(features []Feature) []int64 {
	ids := make([]int64, 0, len(features))
	for _, feature := range features {
		ids = append(ids, feature.ID)
	}
	return ids
}

func collectEntityIDs(entities []Entity) []int64 {
	ids := make([]int64, 0, len(entities))
	for _, entity := range entities {
		ids = append(ids, entity.ID)
	}
	return ids
}

func collectAPIIDs(apis []Api) []int64 {
	ids := make([]int64, 0, len(apis))
	for _, api := range apis {
		ids = append(ids, api.ID)
	}
	return ids
}
