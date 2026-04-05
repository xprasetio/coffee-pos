package handler

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/xprasetio/coffee-pos/config"
	"github.com/xprasetio/coffee-pos/internal/entity"
	"github.com/xprasetio/coffee-pos/internal/middleware"
	"github.com/xprasetio/coffee-pos/internal/repository"
	"github.com/xprasetio/coffee-pos/internal/service"
	"github.com/xprasetio/coffee-pos/pkg/response"
	"github.com/xprasetio/coffee-pos/pkg/validator"
)

// NewRouter creates a new Gin router with configured routes
func NewRouter(db *sql.DB, cfg *config.Config, v *validator.Validator) *gin.Engine {
	// Set gin mode based on environment
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()

	// Add middleware
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, cfg.JWTSecret, cfg.JWTExpiryHours)
	categoryService := service.NewCategoryService(categoryRepo)

	// Initialize handlers
	authHandler := NewAuthHandler(authService, v)
	categoryHandler := NewCategoryHandler(categoryService, v)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		v1.GET("/health", healthHandler)

		// Auth routes (no auth middleware)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Route group untuk owner — semua endpoint di sini butuh login + role owner
		ownerGroup := v1.Group("/owner")
		ownerGroup.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		ownerGroup.Use(middleware.RoleMiddleware(entity.RoleOwner))
		{
			ownerGroup.GET("/categories", categoryHandler.FindAll)
			ownerGroup.GET("/categories/:id", categoryHandler.FindByID)
			ownerGroup.POST("/categories", categoryHandler.Create)
			ownerGroup.PUT("/categories/:id", categoryHandler.Update)
			ownerGroup.DELETE("/categories/:id", categoryHandler.Delete)
		}

		// Route group untuk cashier — semua endpoint di sini butuh login + role cashier
		cashierGroup := v1.Group("/cashier")
		cashierGroup.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		cashierGroup.Use(middleware.RoleMiddleware(entity.RoleCashier))
		_ = cashierGroup
	}

	return router
}

// healthHandler handles health check endpoint
func healthHandler(c *gin.Context) {
	response.OK(c, "server is running", nil)
}
