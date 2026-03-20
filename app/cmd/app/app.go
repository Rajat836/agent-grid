package app

import (
	"app/agent_grid/cmd/app/middlewares"
	"app/agent_grid/internal/clients"
	"app/agent_grid/internal/config"
	"app/agent_grid/internal/controllers"
	"app/agent_grid/internal/repositories"
	"app/agent_grid/internal/services"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"bitbucket.org/fyscal/be-commons/pkg/clickhouse"
	commonClients "bitbucket.org/fyscal/be-commons/pkg/clients"
	"bitbucket.org/fyscal/be-commons/pkg/db"
	"bitbucket.org/fyscal/be-commons/pkg/global"
	"bitbucket.org/fyscal/be-commons/pkg/log"
	networks "bitbucket.org/fyscal/be-commons/pkg/network"
	"bitbucket.org/fyscal/be-commons/pkg/telemetry"
	"google.golang.org/grpc"

	"github.com/gin-gonic/gin"
)

type App struct {
	db           *db.Store
	clickhouseDb *clickhouse.Store
	cache        db.CacheStoreMethods
	fastCache    db.DirtyCacheMethods
	repos        *repositories.Repositories
	networkOps   networks.NetworkOpsMethods

	clients     *clients.Clients
	controllers *controllers.Controllers
	middlewares *middlewares.Middlewares
	services    *services.Services
	router      *gin.Engine
	http        *http.Server
	grpc        *grpc.Server
	telemetry   telemetry.TelemetryUtilsMethods
	sqsClient   commonClients.ClientSqsMethods
}

func (app *App) newDatabaseConnection(cfg *config.Config, logger log.Logger) {
	var err error
	app.db, err = db.NewStore(logger, cfg.Database.MasterDatabaseDsn, cfg.Database.SlaveDatabaseDsn, cfg.AppName)
	if err != nil {
		panic(fmt.Errorf("db initialization failed: %w", err))
	}
}

func (app *App) newCacheConnection(ctx context.Context, cfg *config.Config, logger log.Logger) {
	var err error
	host := cfg.Redis.Host
	port := strconv.Itoa(cfg.Redis.Port)
	username := cfg.Redis.Username
	password := cfg.Redis.Password
	tlsEnabled := cfg.Redis.TlsEnabled
	app.cache, err = db.NewRedisClient(ctx, host, port, username, password, tlsEnabled, logger)
	if err != nil {
		panic(err)
	}
}

func (app *App) newClickHouseConnection(cfg *config.Config, logger log.Logger) {
	var err error
	if !cfg.ClickHouse.Enabled {
		logger.Infof("ClickHouse is disabled, skipping ClickHouse initialization")
		return
	}

	if cfg.ClickHouse.DSN == "" {
		panic(fmt.Errorf("ClickHouse is enabled but DSN not configured, skipping ClickHouse initialization"))
	}

	app.clickhouseDb, err = clickhouse.NewStore(logger, cfg.ClickHouse.DSN, cfg.AppName)
	if err != nil {
		panic(fmt.Errorf("click_house initialization failed: %w", err))
	}

	logger.Infof("ClickHouse connection initialized successfully")
}

func (app *App) setUpHandlers(cfg *config.Config, logger log.Logger) *gin.Engine {
	var router *gin.Engine
	if cfg.Environment == string(global.StageEnv) || cfg.Environment == string(global.ProdEnv) {
		logger.Infof("setting gin to release mode")
		gin.SetMode(gin.ReleaseMode)
	}
	router = gin.Default()
	app.telemetry.EnableGinTracing(router)
	app.addRoutes(router, app.middlewares)
	return router
}

func (app *App) newApp(cfg *config.Config, l log.Logger) {
	ctx := context.Background()
	logger := l.With(ctx)

	app.newDatabaseConnection(cfg, logger)
	app.newClickHouseConnection(cfg, logger)
	app.newSQSClient(cfg, logger)
	app.newCacheConnection(ctx, cfg, logger)
	app.telemetry = telemetry.NewTelemetryUtils(&telemetry.TelemetryUtils{
		Logger:      logger,
		ServiceName: cfg.AppName,
		AppEnv:      cfg.Environment,
		AppVersion:  cfg.AppVersion,
		ExporterURL: cfg.OtlpExporterUrl,
		Ctx:         ctx,
	})

	app.fastCache = db.NewDirtyCache(logger, &app.cache, cfg.AppName)
	app.networkOps = networks.NewNetworkOps(cfg.AppName, logger)
	app.clients = clients.NewClients(cfg, logger, app.cache, app.networkOps)
	app.repos = repositories.NewRepositories(app.db, app.cache, logger, app.fastCache, app.clickhouseDb)
	app.services = services.NewServices(cfg, app.db, app.repos, app.cache, logger, app.clients)
	app.controllers = controllers.NewControllers(cfg, logger, app.services)
	app.middlewares = middlewares.NewMiddlewares(cfg, app.db, app.repos, app.cache, logger, app.clients)
	app.router = app.setUpHandlers(cfg, logger)
	app.http = &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.ServerPort),
		Handler:      app.router,
		ReadTimeout:  time.Second * 60,
		WriteTimeout: time.Second * 60,
		IdleTimeout:  time.Second * 60,
	}
}

func (app *App) start(cfg *config.Config, logger log.Logger) {
	defer func() {
		logger.Infof("app shutting down: cleaning up")
	}()

	logger.Infof("app initialized: HTTP running on port %d, gRPC running on port %d", cfg.ServerPort, cfg.GrpcPort)

	// Start HTTP server (blocking)
	if err := app.http.ListenAndServe(); err != nil {
		logger.Errorf("HTTP server failed to start: %v", err)
	}
}

func (app *App) NewApplication(cfg *config.Config, logger log.Logger) {
	app.newApp(cfg, logger)
	app.start(cfg, logger)

	defer app.destroy(logger)
}

func (app *App) newSQSClient(cfg *config.Config, logger log.Logger) {

	ctx := context.Background()

	logger.Infof("Will setup SQS/SNS client")

	app.sqsClient = commonClients.NewClientSqs(&commonClients.SqsConfig{
		Ctx:     ctx,
		Region:  "",
		Logger:  logger,
		EnvName: global.Environment(cfg.Environment),
	})
}

func (app *App) destroy(logger log.Logger) {
	logger.Infof("shutting down app")
}
