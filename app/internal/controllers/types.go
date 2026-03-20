package controllers

import (
	"app/agent_grid/internal/config"
	"app/agent_grid/internal/services"

	"bitbucket.org/fyscal/be-commons/pkg/log"
)

const (
	QueryParamLimit = "limit"

	PathParamFeatureID = "feature_id"
	PathParamEntityID  = "entity_id"
	PathParamApiID     = "api_id"

	MaxPaginationLimit = 100
)

type ControllerAccess struct {
	Cfg      *config.Config
	Logger   log.Logger
	Services *services.Services
}
