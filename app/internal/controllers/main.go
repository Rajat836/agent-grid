package controllers

import (
	"app/agent_grid/internal/config"
	"app/agent_grid/internal/services"

	"bitbucket.org/fyscal/be-commons/pkg/log"
)

type Controllers struct {
	Health ControllerHealthMethods
}

func NewControllers(cfg *config.Config, logger log.Logger, services *services.Services) *Controllers {
	access := &ControllerAccess{
		Cfg:      cfg,
		Logger:   logger,
		Services: services,
	}
	return &Controllers{
		Health: NewHealthController(access),
	}
}
