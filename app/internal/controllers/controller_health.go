package controllers

import (
	globaltypes "app/agent_grid/internal/global_types"
	"app/agent_grid/internal/response"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// ControllerHealthMethods defines the interface for health controller
type ControllerHealthMethods interface {
	Health(c *gin.Context)
	Readiness(c *gin.Context)
	OntologySummary(c *gin.Context)
}

// healthController implements ControllerHealthMethods
type healthController struct {
	Access *ControllerAccess
}

// NewHealthController creates a new health controller
func NewHealthController(access *ControllerAccess) ControllerHealthMethods {
	return &healthController{
		Access: access,
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
	health := h.Access.Services.Health.GetHealth()
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
	readiness := h.Access.Services.Health.GetReadiness()
	c.JSON(http.StatusOK, readiness)
}

func (h *healthController) OntologySummary(c *gin.Context) {
	log := h.Access.Logger

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  2048,
		WriteBufferSize: 2048,
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for simplicity, adjust for production
		},
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Infof("Failed to upgrade to websocket: %v", err)
		return
	}
	defer conn.Close()

	// Loop to read messages from client
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Infof("Error reading WebSocket message: %v", err)
			break // Exit loop on error (e.g., client disconnected)
		}

		if messageType == websocket.TextMessage {
			var req globaltypes.RequestOntologySummary
			if err := json.Unmarshal(p, &req); err != nil {
				response.SendWebSocketMessage(conn, "error", "Invalid JSON format for prompt", nil)
				continue // Continue listening for next message
			}

			if req.Prompt == "" {
				response.SendWebSocketMessage(conn, "error", "Prompt cannot be empty", nil)
				continue
			}

			// Process the request in a separate goroutine to avoid blocking the read loop
			go h.Access.Services.Health.GetOntologySummary(conn, &req)

		} else {
			log.Infof("Received non-text WebSocket message type: %d", messageType)
			// Optionally send an error or ignore
		}
	}
}
