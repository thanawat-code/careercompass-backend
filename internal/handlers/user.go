package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thanawat-code/careercompass-backend/internal/database"
	"github.com/thanawat-code/careercompass-backend/internal/models"
)

func GetUsers(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()

		query := `SELECT id, email, display_name, gender, created_at, updated_at FROM users ORDER BY created_at DESC`
		rows, err := db.Pool.Query(ctx, query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
			return
		}
		defer rows.Close()

		var users []models.User
		for rows.Next() {
			var user models.User
			if err := rows.Scan(&user.ID, &user.Email, &user.DisplayName, &user.Gender, &user.CreatedAt, &user.UpdatedAt); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan user"})
				return
			}
			users = append(users, user)
		}

		if users == nil {
			users = []models.User{}
		}

		c.JSON(http.StatusOK, users)
	}
}
