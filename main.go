package main

import (
	"log"
	"godoor/config"
	"godoor/models"
	"godoor/routes"
)

func createTestUser() {
	var count int64
	models.DB.Model(&models.User{}).Count(&count)

	if count == 0 {
		log.Println("Creating test user...")
		user, err := models.CreateUser("admin", "admin123", "Admin", "User")
		if err != nil {
			log.Fatal("Failed to create test user:", err)
		}

		// Create user preferences if they don't exist
		_, err = models.GetOrCreateUserPreference(user)
		if err != nil {
			log.Fatal("Failed to create user preferences:", err)
		}

		log.Printf("Test user created: %s (ID: %d)", user.Username, user.ID)

		// Create test space if none exist
		var spaceCount int64
		models.DB.Model(&models.Space{}).Count(&spaceCount)
		if spaceCount == 0 {
			log.Println("Creating test space...")
			space, err := models.CreateSpace(user, "Default", "")
			if err != nil {
				log.Fatal("Failed to create test space:", err)
			}

			// Create test foods
			foods := []struct {
				name        string
				description string
			}{
				{"Tomato", "Fresh red tomato"},
				{"Onion", "Yellow cooking onion"},
				{"Garlic", "Fresh garlic cloves"},
				{"Olive Oil", "Extra virgin olive oil"},
				{"Salt", "Sea salt"},
				{"Black Pepper", "Ground black pepper"},
			}

			for _, foodData := range foods {
				_, err := models.CreateFood(space, foodData.name, foodData.description)
				if err != nil {
					log.Printf("Failed to create food %s: %v", foodData.name, err)
				}
			}
			log.Println("Test foods created")
		}
	}
}

func main() {
	// Initialize database
	config.InitDB()

	// Auto migrate database schema
	models.AutoMigrate()

	// Create test user and preferences if none exist
	createTestUser()

	// Setup routes
	r := routes.SetupRouter()

	// Start server
	log.Println("Server starting on :8080")
	r.Run(":8080")
}
