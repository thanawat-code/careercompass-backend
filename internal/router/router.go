package router

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/thanawat-code/careercompass-backend/internal/config"
	"github.com/thanawat-code/careercompass-backend/internal/database"
	"github.com/thanawat-code/careercompass-backend/internal/handlers"
	"github.com/thanawat-code/careercompass-backend/internal/middleware"
	"github.com/thanawat-code/careercompass-backend/internal/services"

	_ "github.com/thanawat-code/careercompass-backend/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	careerHandler := handlers.NewCareerHandler(db)

	// Health check endpoint (public)
	router.GET("/health", handlers.HealthCheck(db))

	// Swagger endpoint (public)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	api := router.Group("/api")
	{
		// ── Public routes (no auth required) ─────────────────────────────
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Learning path browsing — public so unauthenticated users can explore
		api.GET("/learning-paths", handlers.GetAllLearningPaths(db))
		api.GET("/learning-path/:career_name", handlers.GetLearningPath(db))

		// AI endpoints — public (no user data involved)
		api.POST("/career-recommend", careerHandler.RecommendCareer)
		api.POST("/quiz/generate", handlers.GenerateQuiz)

		// ── Protected routes (require valid JWT Bearer token) ─────────────
		protected := api.Group("/")
		protected.Use(middleware.JWTAuth(authService))
		{
			// User list (contains email addresses — restricted)
			protected.GET("/users", handlers.GetUsers(db))

			// Learning path progress (user-specific data)
			protected.POST("/learning-path/progress", handlers.UpdateUserProgress(db))
			protected.POST("/learning-path/complete-stage", handlers.CompleteStageAndUnlockNext(db))
			protected.GET("/learning-path/progress/:user_id", handlers.GetUserProgress(db))
			protected.DELETE("/learning-path/:career_name/reset", handlers.ResetLearningPath(db))
		}
	}

	return router
}
