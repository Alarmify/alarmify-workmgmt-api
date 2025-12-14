package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	startTime time.Time
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{
		startTime: time.Now(),
	}
}

type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Uptime    string `json:"uptime"`
	Service   string `json:"service"`
	Version   string `json:"version"`
}

// GetHealth handles the health check endpoint
// @Summary Get health status
// @Description Returns comprehensive health status including uptime and service information
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func (h *HealthHandler) GetHealth(c *gin.Context) {
	uptime := time.Since(h.startTime)

	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Uptime:    uptime.String(),
		Service:   "workmgmt-api",
		Version:   "1.0.0",
	}

	c.JSON(http.StatusOK, response)
}

// GetReadiness handles the readiness check endpoint
// @Summary Get readiness status
// @Description Returns readiness probe status for Kubernetes/container orchestration
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health/ready [get]
func (h *HealthHandler) GetReadiness(c *gin.Context) {
	response := gin.H{
		"status":  "ready",
		"service": "workmgmt-api",
	}

	c.JSON(http.StatusOK, response)
}

// GetLiveness handles the liveness check endpoint
// @Summary Get liveness status
// @Description Returns liveness probe status for Kubernetes/container orchestration
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health/live [get]
func (h *HealthHandler) GetLiveness(c *gin.Context) {
	response := gin.H{
		"status":  "alive",
		"service": "workmgmt-api",
	}

	c.JSON(http.StatusOK, response)
}
