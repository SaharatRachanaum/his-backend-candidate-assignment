package handler

import (
	"net/http"

	"github.com/agnos/hospital-middleware/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HealthHandler handles health check probes
type HealthHandler struct {
	db *gorm.DB
}

// NewHealthHandler creates a new HealthHandler instance
func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// HealthCheck handles health probe checks
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	status := gin.H{
		"status":  "healthy",
		"service": "hospital-middleware-api",
	}

	if h.db != nil {
		sqlDB, err := h.db.DB()
		if err != nil || sqlDB.Ping() != nil {
			status["database"] = "down"
			response.Error(c, http.StatusServiceUnavailable, "Database connection error", status)
			return
		}
		status["database"] = "up"
	}

	response.Success(c, http.StatusOK, "Service is healthy", status)
}
