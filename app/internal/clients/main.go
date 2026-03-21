package clients

import (
	"app/agent_grid/internal/config"
	"context"

	"bitbucket.org/fyscal/be-commons/pkg/clients"
	"bitbucket.org/fyscal/be-commons/pkg/db"
	"bitbucket.org/fyscal/be-commons/pkg/global"
	"bitbucket.org/fyscal/be-commons/pkg/log"
	networks "bitbucket.org/fyscal/be-commons/pkg/network"
)

type Clients struct {
	ClientSqs     clients.ClientSqsMethods
	ClientGeneral ClientGeneralMethods
	ClientAgent   ClientAgentMethods
}

type clientAccess struct {
	cfg        *config.Config
	logger     log.Logger
	cacheStore db.CacheStoreMethods
	networkOps networks.NetworkOpsMethods
}

func NewClients(cfg *config.Config, logger log.Logger, cacheStore db.CacheStoreMethods, ops networks.NetworkOpsMethods) *Clients {
	access := &clientAccess{
		cfg:        cfg,
		logger:     logger,
		cacheStore: cacheStore,
		networkOps: ops,
	}

	ctx := context.TODO()

	sqsClient := clients.NewClientSqs(&clients.SqsConfig{Ctx: ctx, Region: cfg.SQS.Region, Logger: logger, EnvName: global.Environment(cfg.Environment)})

	return &Clients{
		ClientSqs:     sqsClient,
		ClientGeneral: NewClientGeneral(access),
		ClientAgent:   NewClientAgent(access),
	}
}
