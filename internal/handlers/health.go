package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thanawat-code/careercompass-backend/internal/database"
)

// HealthCheck godoc
// @Summary Check API Health
// @Description get the status of server and database connection
// @Tags System
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 503 {object} map[string]interface{}
// @Router /health [get]
func HealthCheck(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()

		// Check database connection
		if err := db.Ping(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":   "error",
				"database": "disconnected",
				"error":    err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":   "ok",
			"database": "connected",
		})
	}
}
