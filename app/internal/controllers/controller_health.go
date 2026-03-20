package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"app/ontology_bot/internal/services"
)

// ControllerHealthMethods defines the interface for health controller
type ControllerHealthMethods interface {
	Health(c *gin.Context)
	Readiness(c *gin.Context)
}

// healthController implements ControllerHealthMethods
type healthController struct {
	healthService services.ServiceHealthMethods
}

// NewHealthController creates a new health controller
func NewHealthController(version string) ControllerHealthMethods {
	healthService := services.NewHealthService(version)
	return &healthController{
		healthService: healthService,
	}
}

// Health returns the health status of the service
// @Summary Health Check
// @Description Returns the health status of the service
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} services.HealthCheckResponse
// @Router /ontology_bot/v1/health [get]
func (h *healthController) Health(c *gin.Context) {
	health := h.healthService.GetHealth()
	c.JSON(http.StatusOK, health)
}

// Readiness returns the readiness status of the service
// @Summary Readiness Check
// @Description Returns whether the service is ready to handle requests
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} services.ReadinessResponse
// @Router /ontology_bot/v1/health/ready [get]
func (h *healthController) Readiness(c *gin.Context) {
	readiness := h.healthService.GetReadiness()
	c.JSON(http.StatusOK, readiness)
}
