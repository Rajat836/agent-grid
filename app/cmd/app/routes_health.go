package app

import (
	"app/ontology_bot/internal/controllers"

	"github.com/gin-gonic/gin"
)

func registerHealthRoutes(engine *gin.Engine, version string) {
	healthController := controllers.NewHealthController(version)

	// Public routes - no authentication required
	public := engine.Group("/ontology_bot/v1")
	{
		public.GET("/health", healthController.Health)
		public.GET("/health/ready", healthController.Readiness)
	}
}
