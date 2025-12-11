package routes

import (
	"godoor/config"
	"godoor/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRouter configures all routes
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Set maximum multipart memory (32MB for image uploads)
	r.MaxMultipartMemory = 32 << 20

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

	// Static files
	r.Static(config.MediaURL, config.MediaRoot)

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
			// API root endpoint - return available endpoints like Django REST Framework
			rootResponse := gin.H{
				"automation":           "/api/automation/",
				"bookmarklet-import":   "/api/bookmarklet-import/",
				"cook-log":             "/api/cook-log/",
				"custom-filter":        "/api/custom-filter/",
				"food":                 "/api/food/",
				"food-inherit-field":   "/api/food-inherit-field/",
				"group":                "/api/group/",
				"import-log":           "/api/import-log/",
				"ingredient":           "/api/ingredient/",
				"invite-link":          "/api/invite-link/",
				"keyword":              "/api/keyword/",
				"meal-plan":            "/api/meal-plan/",
				"meal-type":            "/api/meal-type/",
				"recipe":               "/api/recipe/",
				"recipe-book":          "/api/recipe-book/",
				"recipe-book-entry":    "/api/recipe-book-entry/",
				"server-settings":      "/api/server-settings/",
				"shopping-list":        "/api/shopping-list/",
				"shopping-list-entry":  "/api/shopping-list-entry/",
				"shopping-list-recipe": "/api/shopping-list-recipe/",
				"space":                "/api/space/",
				"step":                 "/api/step/",
				"storage":              "/api/storage/",
				"supermarket":          "/api/supermarket/",
				"supermarket-category": "/api/supermarket-category/",
				"unit":                 "/api/unit/",
				"unit-conversion":      "/api/unit-conversion/",
				"user":                 "/api/user/",
				"user-preference":      "/api/user-preference/",
				"user-space":           "/api/user-space/",
			}
			c.JSON(200, rootResponse)
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

		// Keyword routes
		protected.GET("/keyword/", GetKeywords)
		protected.POST("/keyword/", CreateKeyword)
		protected.GET("/keyword/:id/", GetKeyword)
		protected.PUT("/keyword/:id/", UpdateKeyword)
		protected.PATCH("/keyword/:id/", UpdateKeyword)
		protected.DELETE("/keyword/:id/", DeleteKeyword)

		// Cook log routes
		protected.GET("/cook-log/", GetCookLogs)
		protected.POST("/cook-log/", CreateCookLog)
		protected.GET("/cook-log/:id/", GetCookLog)
		protected.PUT("/cook-log/:id/", UpdateCookLog)
		protected.PATCH("/cook-log/:id/", UpdateCookLog)
		protected.DELETE("/cook-log/:id/", DeleteCookLog)

		// Meal type routes
		protected.GET("/meal-type/", GetMealTypes)
		protected.POST("/meal-type/", CreateMealType)

		// Meal plan routes
		protected.GET("/meal-plan/", GetMealPlans)
		protected.POST("/meal-plan/", CreateMealPlan)
		protected.GET("/meal-plan/:id/", GetMealPlan)
		protected.PUT("/meal-plan/:id/", UpdateMealPlan)
		protected.PATCH("/meal-plan/:id/", UpdateMealPlan)
		protected.DELETE("/meal-plan/:id/", DeleteMealPlan)

		// Step routes
		protected.GET("/step/", GetSteps)
		protected.POST("/step/", CreateStep)
		protected.GET("/step/:id/", GetStep)
		protected.PUT("/step/:id/", UpdateStep)
		protected.PATCH("/step/:id/", UpdateStep)
		protected.DELETE("/step/:id/", DeleteStep)

		// Food routes
		protected.GET("/food/", GetFoods)
		protected.POST("/food/", CreateFood)
		protected.GET("/food/:id/", GetFood)
		protected.PUT("/food/:id/", UpdateFood)
		protected.PATCH("/food/:id/", UpdateFood)
		protected.DELETE("/food/:id/", DeleteFood)

		// Recipe book entry routes
		protected.GET("/recipe-book-entry/", GetRecipeBookEntries)
		protected.POST("/recipe-book-entry/", CreateRecipeBookEntry)
		protected.DELETE("/recipe-book-entry/:id/", DeleteRecipeBookEntry)

		// Recipe routes
		protected.GET("/recipe/", GetRecipes)
		protected.POST("/recipe/", CreateRecipe)
		protected.GET("/recipe/:id/", GetRecipe)
		protected.PUT("/recipe/:id/", UpdateRecipe)
		protected.PATCH("/recipe/:id/", UpdateRecipe)
		protected.DELETE("/recipe/:id/", DeleteRecipe)
		protected.PUT("/recipe/:id/shopping/", RecipeShopping)
		protected.PUT("/recipe/:id/image/", RecipeImage)
	}

	return r
}
