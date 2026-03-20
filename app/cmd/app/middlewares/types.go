package middlewares

import (
	"app/agent_grid/internal/clients"
	"app/agent_grid/internal/config"
	"app/agent_grid/internal/repositories"

	"bitbucket.org/fyscal/be-commons/pkg/db"
	"bitbucket.org/fyscal/be-commons/pkg/log"
)

type MiddlewareAccess struct {
	Cfg          *config.Config
	Db           *db.Store
	Cache        db.CacheStoreMethods
	Logger       log.Logger
	Clients      *clients.Clients
	Repositories repositories.Repositories
}
