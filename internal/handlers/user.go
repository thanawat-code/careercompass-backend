package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thanawat-code/careercompass-backend/internal/database"
	"github.com/thanawat-code/careercompass-backend/internal/models"
)

func GetUsers(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()

		query := `SELECT id, name, email, created_at, updated_at FROM users ORDER BY created_at DESC`
		rows, err := db.Pool.Query(ctx, query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
			return
		}
		defer rows.Close()

		var users []models.User
		for rows.Next() {
			var user models.User
			if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
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

func CreateUser(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.CreateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx := context.Background()
		user := models.User{
			ID:        uuid.New(),
			Name:      req.Name,
			Email:     req.Email,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		query := `INSERT INTO users (id, name, email, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`
		_, err := db.Pool.Exec(ctx, query, user.ID, user.Name, user.Email, user.CreatedAt, user.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		c.JSON(http.StatusCreated, user)
	}
}
