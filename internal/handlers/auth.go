package handlers

import (
	"context"
	"log"
	"net/http"
	"strings"
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

// Register godoc
// @Summary User Registration
// @Description Register a new user
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Registration Info"
// @Success 201 {object} models.AuthResponse
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Normalize email to lowercase (fixes REG_008: case-insensitive duplicate check)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

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
		log.Printf("❌ Register - DB check error: %v", err)
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
		log.Printf("❌ Register - DB insert error: %v", err)
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

// Login godoc
// @Summary User Login
// @Description Authenticate user and return JWT
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login Credentials"
// @Success 200 {object} models.AuthResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Normalize email to lowercase (fixes LOG_008: case-insensitive login)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

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
