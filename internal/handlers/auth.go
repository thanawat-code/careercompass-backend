package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/thanawat-code/careercompass-backend/internal/database"
	"github.com/thanawat-code/careercompass-backend/internal/models"
	"github.com/thanawat-code/careercompass-backend/internal/services"
)

type AuthHandler struct {
	db          *database.DB
	authService *services.AuthService
}

func NewAuthHandler(db *database.DB, authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		db:          db,
		authService: authService,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate password match
	if err := h.authService.ValidatePasswordMatch(req.Password, req.ConfirmPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "passwords do not match"})
		return
	}

	// Check if user already exists
	ctx := context.Background()
	var existingID uuid.UUID
	err := h.db.Pool.QueryRow(ctx, "SELECT id FROM users WHERE email = $1", req.Email).Scan(&existingID)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
		return
	} else if err != pgx.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check existing user"})
		return
	}

	// Hash password
	hashedPassword, err := h.authService.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	// Create user
	user := models.User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: hashedPassword,
		DisplayName:  req.DisplayName,
		Gender:       req.Gender,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	query := `INSERT INTO users (id, email, password_hash, display_name, gender, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err = h.db.Pool.Exec(ctx, query, user.ID, user.Email, user.PasswordHash, user.DisplayName, user.Gender, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	// Generate JWT token
	token, err := h.authService.GenerateToken(user.ID.String(), user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, models.AuthResponse{
		User:  user,
		Token: token,
	})
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user by email
	ctx := context.Background()
	var user models.User
	query := `SELECT id, email, password_hash, display_name, gender, created_at, updated_at 
	          FROM users WHERE email = $1`
	err := h.db.Pool.QueryRow(ctx, query, req.Email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.DisplayName, &user.Gender, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find user"})
		return
	}

	// Verify password
	if err := h.authService.VerifyPassword(user.PasswordHash, req.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Generate JWT token
	token, err := h.authService.GenerateToken(user.ID.String(), user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, models.AuthResponse{
		User:  user,
		Token: token,
	})
}
