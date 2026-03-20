package repositories

import (
	"bitbucket.org/fyscal/be-commons/pkg/clickhouse"
	"bitbucket.org/fyscal/be-commons/pkg/db"
	"bitbucket.org/fyscal/be-commons/pkg/log"
)

type RepositoryAccess struct {
	Db           *db.Store
	Cache        db.CacheStoreMethods
	Logger       log.Logger
	FastCache    db.DirtyCacheMethods
	ClickHouseDb *clickhouse.Store
}
