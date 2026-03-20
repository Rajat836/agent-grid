package services

import (
	"app/agent_grid/internal/clients"
	"app/agent_grid/internal/config"
	"app/agent_grid/internal/repositories"
	"time"

	"bitbucket.org/fyscal/be-commons/pkg/db"
	"bitbucket.org/fyscal/be-commons/pkg/log"
)

type ServiceAccess struct {
	Cfg          *config.Config
	Db           *db.Store
	Cache        db.CacheStoreMethods
	Logger       log.Logger
	Clients      *clients.Clients
	Repositories repositories.Repositories
}

const (
	ParamName         = "name"
	ParamDescription  = "description"
	ParamCode         = "code"
	ParamDisplayOrder = "display_order"
	ParamIsStart      = "is_start"
	ParamIsTerminal   = "is_terminal"
	ParamFeatureID    = "feature_id"
)

type HealthCheckResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
}

// ReadinessResponse contains readiness information
type ReadinessResponse struct {
	Ready     bool      `json:"ready"`
	Timestamp time.Time `json:"timestamp"`
}
