package router

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/thanawat-code/careercompass-backend/internal/config"
	"github.com/thanawat-code/careercompass-backend/internal/database"
	"github.com/thanawat-code/careercompass-backend/internal/handlers"
	"github.com/thanawat-code/careercompass-backend/internal/services"
)

func Setup(cfg *config.Config, db *database.DB) *gin.Engine {
	// Set Gin mode
	gin.SetMode(cfg.Server.GinMode)

	router := gin.Default()

	// CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORS.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Initialize auth service
	authService, err := services.NewAuthService(cfg.JWT.Secret, cfg.JWT.Expiration)
	if err != nil {
		log.Fatalf("Failed to create auth service: %v", err)
	}

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db, authService)

	// Health check endpoint
	router.GET("/health", handlers.HealthCheck(db))
	careerHandler := handlers.NewCareerHandler(db)

	// API routes
	api := router.Group("/api")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// User routes
		api.GET("/users", handlers.GetUsers(db))
		api.POST("/career-recommend", careerHandler.RecommendCareer)

		// Learning Path routes
		api.GET("/learning-paths", handlers.GetAllLearningPaths(db))
		api.GET("/learning-path/:career_name", handlers.GetLearningPath(db))
		api.POST("/learning-path/progress", handlers.UpdateUserProgress(db))
		api.POST("/learning-path/complete-stage", handlers.CompleteStageAndUnlockNext(db))
		api.GET("/learning-path/progress/:user_id", handlers.GetUserProgress(db))
	}

	return router
}
