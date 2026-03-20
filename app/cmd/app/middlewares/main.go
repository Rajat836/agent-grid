package middlewares

import (
	"app/agent_grid/internal/clients"
	"app/agent_grid/internal/config"
	"app/agent_grid/internal/repositories"

	"bitbucket.org/fyscal/be-commons/pkg/db"
	"bitbucket.org/fyscal/be-commons/pkg/log"
)

type Middlewares struct {
}

func NewMiddlewares(cfg *config.Config, db *db.Store, r *repositories.Repositories, cs db.CacheStoreMethods, l log.Logger, c *clients.Clients) *Middlewares {
	_ = &MiddlewareAccess{
		Cfg:          cfg,
		Db:           db,
		Cache:        cs,
		Logger:       l,
		Repositories: *r,
		Clients:      c,
	}
	return &Middlewares{}
}
