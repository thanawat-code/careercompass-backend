package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thanawat-code/careercompass-backend/internal/database"
)

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
