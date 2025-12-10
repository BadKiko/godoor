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
