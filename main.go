package main

import (
	"log"
	"os"
	"godoor/config"
	"godoor/models"
	"godoor/routes"
	"github.com/joho/godotenv"
)

func initializeDatabase() {
	// Check if we need to create admin user
	var userCount int64
	models.DB.Model(&models.User{}).Count(&userCount)

	if userCount == 0 {
		// Create admin user from environment variables or defaults
		adminUsername := os.Getenv("ADMIN_USERNAME")
		adminPassword := os.Getenv("ADMIN_PASSWORD")

		if adminUsername == "" {
			adminUsername = "admin"
		}
		if adminPassword == "" {
			adminPassword = "admin123"
		}

		log.Printf("Creating admin user: %s", adminUsername)
		user, err := models.CreateUser(adminUsername, adminPassword, "Admin", "User")
		if err != nil {
			log.Printf("Failed to create admin user: %v", err)
			return
		}

		// Create user preferences
		_, err = models.GetOrCreateUserPreference(user)
		if err != nil {
			log.Printf("Failed to create user preferences: %v", err)
		}

		// Create default space for admin
		userSpace, err := models.CreateSpaceForUser(user, nil)
		if err != nil {
			log.Printf("Failed to create default space: %v", err)
			return
		}

		log.Printf("Admin user and default space created successfully")
		log.Printf("Login: %s, Password: %s", adminUsername, adminPassword)

		// Create default meal types for the new space
		createDefaultMealTypes(userSpace.Space)
	} else {
		// Check if we need to create default meal types for existing spaces
		var spaceCount int64
		models.DB.Model(&models.Space{}).Count(&spaceCount)

		if spaceCount > 0 {
			// Get first space
			var space models.Space
			if err := models.DB.First(&space).Error; err != nil {
				log.Printf("Failed to get space: %v", err)
				return
			}

			createDefaultMealTypes(space)
		}
	}
}

func createDefaultMealTypes(space models.Space) {
	// Check if we need to create default meal types
	var mealTypeCount int64
	models.DB.Model(&models.MealType{}).Where("space_id = ?", space.ID).Count(&mealTypeCount)

	if mealTypeCount == 0 {
		log.Println("Creating default meal types...")

		// Get first user for created_by
		var user models.User
		if err := models.DB.First(&user).Error; err != nil {
			log.Printf("Failed to get user for meal types: %v", err)
			return
		}

		// Create default meal types
		defaultMealTypes := []struct {
			name  string
			order int
			time  *string
			color *string
			def   bool
		}{
			{"Breakfast", 1, stringPtr("08:00:00"), stringPtr("#FF6B35"), true},
			{"Lunch", 2, stringPtr("12:00:00"), stringPtr("#F7931E"), false},
			{"Dinner", 3, stringPtr("18:00:00"), stringPtr("#FFD23F"), false},
			{"Snack", 4, nil, stringPtr("#06FFA5"), false},
		}

		for _, mt := range defaultMealTypes {
			mealType := models.MealType{
				Name:      mt.name,
				Order:     mt.order,
				Time:      mt.time,
				Color:     mt.color,
				Default:   mt.def,
				CreatedBy: user,
				Space:     space,
			}

			if err := models.DB.Create(&mealType).Error; err != nil {
				log.Printf("Failed to create meal type %s: %v", mt.name, err)
				continue
			}
			log.Printf("Created meal type: %s", mt.name)
		}
	}
}

// stringPtr creates a string pointer
func stringPtr(s string) *string {
	return &s
}

func main() {
	// Load .env file if exists
	if err := godotenv.Load(); err != nil {
		// .env file not found, continue without it
		log.Println("No .env file found, using default environment variables")
	}

	// Initialize database
	config.InitDB()

	// Auto migrate database schema
	models.AutoMigrate()

	// Initialize database with default data
	initializeDatabase()

	// Setup routes
	r := routes.SetupRouter()

	// Get port from environment or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("Server starting on :%s", port)
	r.Run(":" + port)
}
