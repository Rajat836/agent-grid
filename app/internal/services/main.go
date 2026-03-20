package services

import (
	"app/agent_grid/internal/clients"
	"app/agent_grid/internal/config"
	"app/agent_grid/internal/repositories"

	"bitbucket.org/fyscal/be-commons/pkg/db"
	"bitbucket.org/fyscal/be-commons/pkg/log"
)

type Services struct {
	Health ServiceHealthMethods
}

func NewServices(cfg *config.Config, db *db.Store, r *repositories.Repositories, cs db.CacheStoreMethods, l log.Logger, c *clients.Clients) *Services {
	access := &ServiceAccess{
		Cfg:          cfg,
		Db:           db,
		Cache:        cs,
		Logger:       l,
		Repositories: *r,
		Clients:      c,
	}

	return &Services{
		Health: NewServiceHealth(access),
	}
}
