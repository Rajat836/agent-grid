package app

import (
	"app/agent_grid/cmd/app/middlewares"

	"github.com/gin-gonic/gin"
)

func (app *App) addBaseRoutes(router *gin.Engine, middlewares *middlewares.Middlewares) {
	// Public routes - no authentication required
	public := router.Group("/agent/v1")
	{
		public.GET("/health", app.controllers.Health.Health)
		public.GET("/health/ready", app.controllers.Health.Readiness)
		public.GET("/ontology/summary", app.controllers.Health.OntologySummary)
	}
}
