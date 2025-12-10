package routes

import (
	"net/http"
	"godoor/middleware"
	"github.com/gin-gonic/gin"
)

// SetupRouter configures all routes
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Authentication routes (no auth required)
	r.POST("/api-token-auth/", AuthenticateUser)
	r.Any("/api-auth/", func(c *gin.Context) {
		// DRF auth endpoints - simplified response
		c.JSON(http.StatusOK, gin.H{"message": "Django REST framework authentication"})
	})

	// Public API routes (no auth required)
	api := r.Group("/api")
	{
		api.GET("/server-settings/current/", ServerSettingsCurrent)
		api.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusForbidden, gin.H{"detail": "Authentication credentials were not provided."})
		})
	}

	// Protected API routes with space scoping
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		// User routes
		protected.GET("/user/", GetUsers)
		protected.GET("/user/:id/", GetUser)
		protected.PATCH("/user/:id/", UpdateUser)

		// Space routes
		protected.GET("/space/", GetSpaces)
		protected.POST("/space/", CreateSpace)
		protected.GET("/space/:id/", GetSpace)
		protected.PUT("/space/:id/", UpdateSpace)
		protected.PATCH("/space/:id/", UpdateSpace)
		protected.GET("/space/current/", GetCurrentSpace)

		// UserSpace routes
		protected.GET("/user-space/", GetUserSpaces)
		protected.GET("/user-space/:id/", GetUserSpace)
		protected.PUT("/user-space/:id/", UpdateUserSpace)
		protected.PATCH("/user-space/:id/", UpdateUserSpace)
	}

	// Protected API routes without space scoping (like user preferences)
	protectedNoScope := r.Group("/api")
	protectedNoScope.Use(middleware.AuthMiddleware())
	{
		// User preference routes (no scope middleware needed)
		protectedNoScope.GET("/user-preference/", GetUserPreference)
		protectedNoScope.PATCH("/user-preference/", UpdateUserPreference)

		// Recipe book routes
		protected.GET("/recipe-book/", GetRecipeBooks)
		protected.POST("/recipe-book/", CreateRecipeBook)
		protected.GET("/recipe-book/:id/", GetRecipeBook)
		protected.PUT("/recipe-book/:id/", UpdateRecipeBook)
		protected.PATCH("/recipe-book/:id/", UpdateRecipeBook)
		protected.DELETE("/recipe-book/:id/", DeleteRecipeBook)
	}

	return r
}
