package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/xprasetio/coffee-pos/pkg/response"
)

// NewRouter creates a new Gin router with configured routes
func NewRouter(appEnv string) *gin.Engine {
	// Set gin mode based on environment
	if appEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()

	// Add middleware
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		v1.GET("/health", healthHandler)
	}

	return router
}

// healthHandler handles health check endpoint
func healthHandler(c *gin.Context) {
	response.OK(c, "server is running", nil)
}
