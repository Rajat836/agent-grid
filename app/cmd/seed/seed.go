package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
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
	TransitionType       string    `gorm:"type:varchar(64);index:idx_entity_transitions_transition_type" json:"transition_type"`
	CreatedAt            time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	FromEntity           Entity    `gorm:"foreignKey:FromEntityID" json:"from_entity,omitempty"`
	ToEntity             Entity    `gorm:"foreignKey:ToEntityID" json:"to_entity,omitempty"`
}

func (EntityTransition) TableName() string { return "entity_transitions" }

type Entity struct {
	ID           int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	FeatureID    *int64         `gorm:"default:NULL;index:idx_entities_feature_id" json:"feature_id"`
	Code         string         `gorm:"type:varchar(32);not null;uniqueIndex" json:"code"`
	Name         string         `gorm:"type:varchar(64);not null;index:idx_entities_name" json:"name"`
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

func (FeatureTeamRole) TableName() string { return "feature_team_roles" }

type Feature struct {
	ID          int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Code        string         `gorm:"type:varchar(32);not null;uniqueIndex" json:"code"`
	Name        string         `gorm:"type:varchar(64);not null;index:idx_features_name" json:"name"`
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
	RelationType     string     `gorm:"type:varchar(64);not null;index:idx_kpis_relationships_relation_type" json:"relation_type"`
	TargetType       string     `gorm:"type:varchar(64);not null;index:idx_kpis_relationships_target" json:"target_type"`
	TargetID         int64      `gorm:"not null;index:idx_kpis_relationships_target" json:"target_id"`
	Weight           float64    `gorm:"default:0" json:"weight"`
	WeightSetBy      string     `gorm:"type:varchar(64)" json:"weight_set_by"`
	WeightReviewedBy string     `gorm:"type:varchar(64)" json:"weight_reviewed_by"`
	WeightReviewedAt *time.Time `gorm:"default:NULL" json:"weight_reviewed_at,omitempty"`
	CreatedAt        time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	Kpi              Kpi        `gorm:"foreignKey:KpiID" json:"kpi,omitempty"`
}

func (KpiRelationship) TableName() string { return "kpis_relationships" }

type Kpi struct {
	ID               int64             `gorm:"primaryKey;autoIncrement" json:"id"`
	Code             string            `gorm:"type:varchar(32);not null;uniqueIndex" json:"code"`
	Name             string            `gorm:"type:varchar(64);not null;index:idx_kpis_name" json:"name"`
	Description      string            `gorm:"type:text" json:"description"`
	MetricType       string            `gorm:"type:varchar(64);not null;index:idx_kpis_metric_type" json:"metric_type"`
	Unit             string            `gorm:"type:varchar(32)" json:"unit"`
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
	Role       string     `gorm:"type:varchar(64);not null;index:idx_service_team_roles_role" json:"role"`
	AssignedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"assigned_at"`
	RevokedAt  *time.Time `gorm:"default:NULL" json:"revoked_at,omitempty"`
	Service    Service    `gorm:"foreignKey:ServiceID" json:"service,omitempty"`
	Team       Team       `gorm:"foreignKey:TeamID" json:"team,omitempty"`
}

func (ServiceTeamRole) TableName() string { return "service_team_roles" }

type Service struct {
	ID               int64             `gorm:"primaryKey;autoIncrement" json:"id"`
	Code             string            `gorm:"type:varchar(32);not null;uniqueIndex" json:"code"`
	Name             string            `gorm:"type:varchar(64);not null;index:idx_services_name" json:"name"`
	Description      string            `gorm:"type:text" json:"description"`
	RepositoryURL    string            `gorm:"type:varchar(512)" json:"repository_url"`
	CriticalityLevel string            `gorm:"type:varchar(32);index:idx_services_criticality_level" json:"criticality_level"`
	CreatedAt        time.Time         `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        time.Time         `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt        gorm.DeletedAt    `gorm:"default:NULL" json:"deleted_at"`
	ServiceTeamRoles []ServiceTeamRole `gorm:"foreignKey:ServiceID" json:"service_team_roles,omitempty"`
}

func (Service) TableName() string { return "services" }

type Team struct {
	ID               int64             `gorm:"primaryKey;autoIncrement" json:"id"`
	Name             string            `gorm:"type:varchar(64);not null;index:idx_teams_name" json:"name"`
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
	numTeams              = 50
	numFeatures           = 50
	numServices           = 50
	numApis               = 50
	numKpis               = 50
	numEntities           = 60
	numServiceDeployments = 50
	numApiDeployments     = 50
	numApiExecutions      = 100
	numApiDbInteractions  = 150
	numKpiRelationships   = 50
	numServiceTeamRoles   = 50
	// numFeatureTeamRoles   = 50
	numEntityApis        = 50
	numEntityTransitions = 50
	numFeatureInstances  = 50
	numEntityInstances   = 50
	numApiMetrics        = 50
	numEntityMetrics     = 50
	numFeatureMetrics    = 50
)

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
		if err := seedKpiRelationships(pgDB, kpis); err != nil {
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
		if err := seedEntityApis(pgDB, entities, apis); err != nil {
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
		list = append(list, Team{Name: fmt.Sprintf("Team-%d", i), SlackChannel: "#alerts", PagerdutyKey: "pk_123"})
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
		list = append(list, Feature{Code: fmt.Sprintf("FEAT_%d", i), Name: "Feature Name", IsActive: true})
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
		list = append(list, Service{Code: fmt.Sprintf("SVC_%d", i), Name: "Service Name", CriticalityLevel: "P0"})
	}
	if err := db.Create(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func seedApis(db *gorm.DB) ([]Api, error) {
	fmt.Println("Seeding apis...")
	var list []Api
	methods := []string{"GET", "POST", "PUT"}
	for i := 1; i <= numApis; i++ {
		list = append(list, Api{Endpoint: fmt.Sprintf("/v1/resource/%d", i), HttpMethod: methods[rand.Intn(len(methods))]})
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
		f := features[rand.Intn(len(features))]
		list = append(list, Entity{FeatureID: &f.ID, Code: fmt.Sprintf("ENT_%d", i), Name: "Entity Name"})
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
		s := svcs[rand.Intn(len(svcs))]
		list = append(list, ServiceDeployment{ID: int64(i), ServiceID: s.ID, Environment: "prod", DeployedAt: time.Now().Unix()})
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
		list = append(list, ApisDeployment{
			ID:                  uuid.NewString(),
			ApiID:               apis[rand.Intn(len(apis))].ID,
			ServiceDeploymentID: deploys[rand.Intn(len(deploys))].ID,
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
	for i := 1; i <= numApiExecutions; i++ {
		list = append(list, ApiExecution{
			ID:               uuid.NewString(),
			EntityInstanceID: ei[rand.Intn(len(ei))].ID,
			ApiDeploymentsID: ad[rand.Intn(len(ad))].ID,
			HttpStatus:       200,
			RequestTimestamp: time.Now().Unix(),
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
	for i := 1; i <= numApiDbInteractions; i++ {
		list = append(list, ApiDbInteraction{
			ID:             uuid.NewString(),
			ApiExecutionID: execs[rand.Intn(len(execs))].ID,
			DatabaseType:   "PostgreSQL",
			DbTableName:    "users",
			OperationType:  "SELECT",
			CreatedAt:      time.Now().Unix(),
		})
	}
	return db.Create(&list).Error
}

func seedKpis(db *gorm.DB) ([]Kpi, error) {
	fmt.Println("Seeding kpis...")
	var list []Kpi
	for i := 1; i <= numKpis; i++ {
		list = append(list, Kpi{Code: fmt.Sprintf("KPI_%d", i), Name: "KPI Name", MetricType: "Latency"})
	}
	if err := db.Create(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func seedKpiRelationships(db *gorm.DB, kpis []Kpi) error {
	fmt.Println("Seeding kpi relationships...")
	var list []KpiRelationship
	for i := 1; i <= numKpiRelationships; i++ {
		list = append(list, KpiRelationship{
			KpiID:        kpis[rand.Intn(len(kpis))].ID,
			RelationType: "impact",
			TargetType:   "service",
			TargetID:     int64(i),
		})
	}
	return db.Create(&list).Error
}

func seedServiceTeamRoles(db *gorm.DB, svcs []Service, teams []Team) error {
	fmt.Println("Seeding service team roles...")
	var list []ServiceTeamRole
	for i := 1; i <= numServiceTeamRoles; i++ {
		list = append(list, ServiceTeamRole{
			ServiceID: svcs[rand.Intn(len(svcs))].ID,
			TeamID:    teams[rand.Intn(len(teams))].ID,
			Role:      "Owner",
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

func seedEntityApis(db *gorm.DB, ents []Entity, apis []Api) error {
	fmt.Println("Seeding entity apis...")
	var list []EntityApi
	for i := 1; i <= numEntityApis; i++ {
		list = append(list, EntityApi{
			EntityID: ents[rand.Intn(len(ents))].ID,
			ApiID:    apis[rand.Intn(len(apis))].ID,
		})
	}
	return db.Create(&list).Error
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
	for i := 1; i <= numEntityTransitions; i++ {
		list = append(list, EntityTransition{
			FromEntityID:   ents[rand.Intn(len(ents))].ID,
			ToEntityID:     ents[rand.Intn(len(ents))].ID,
			TransitionType: "AUTOMATIC",
		})
	}
	return db.Create(&list).Error
}

func seedFeatureInstances(db *gorm.DB, feats []Feature) ([]FeatureInstance, error) {
	fmt.Println("Seeding feature instances...")
	var list []FeatureInstance
	for i := 1; i <= numFeatureInstances; i++ {
		list = append(list, FeatureInstance{
			ID:        uuid.NewString(),
			FeatureID: feats[rand.Intn(len(feats))].ID,
			Status:    "ACTIVE",
			StartedAt: time.Now().Unix(),
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
	for i := 1; i <= numEntityInstances; i++ {
		list = append(list, EntityInstance{
			ID:                uuid.NewString(),
			FeatureInstanceID: fi[rand.Intn(len(fi))].ID,
			EntityID:          ents[rand.Intn(len(ents))].ID,
			Status:            "COMPLETED",
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
		list = append(list, ApiMetric{
			ApisDeploymentID: ad[rand.Intn(len(ad))].ID,
			WindowStart:      time.Now(),
			TotalCalls:       1000,
			SuccessRate:      0.99,
		})
	}
	return db.Create(&list).Error
}

func seedEntityMetrics(db *gorm.DB, ents []Entity) error {
	fmt.Println("Seeding entity metrics...")
	var list []EntityMetric
	for i := 1; i <= numEntityMetrics; i++ {
		list = append(list, EntityMetric{
			EntityID:    ents[rand.Intn(len(ents))].ID,
			WindowStart: time.Now(),
			TotalCount:  500,
			SuccessRate: 0.95,
		})
	}
	return db.Create(&list).Error
}

func seedFeatureMetrics(db *gorm.DB, feats []Feature) error {
	fmt.Println("Seeding feature metrics...")
	var list []FeatureMetric
	for i := 1; i <= numFeatureMetrics; i++ {
		list = append(list, FeatureMetric{
			FeatureID:   feats[rand.Intn(len(feats))].ID,
			WindowStart: time.Now(),
			TotalCount:  200,
			SuccessRate: 0.98,
		})
	}
	return db.Create(&list).Error
}
