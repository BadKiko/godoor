package main

import (
	"log"
	"godoor/config"
	"godoor/models"
	"godoor/routes"
)

// stringPtr creates a string pointer
func stringPtr(s string) *string {
	return &s
}

func createTestUser() {
	var count int64
	models.DB.Model(&models.User{}).Count(&count)
	log.Printf("User count: %d", count)

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

	// Always try to create test space and data
	var spaceCount int64
	models.DB.Model(&models.Space{}).Count(&spaceCount)
	log.Printf("Space count: %d", spaceCount)

	// Get the user (should exist now)
	var user models.User
	if err := models.DB.First(&user).Error; err != nil {
		log.Printf("Failed to get user: %v", err)
		return
	}

	var space *models.Space
	if spaceCount == 0 {
		log.Println("Creating test space...")
		var err error
		userSpace, err := models.CreateSpaceForUser(&user, nil)
		if err != nil {
			log.Printf("Failed to create test space: %v", err)
			return
		}
		space = &userSpace.Space
		log.Printf("Test space created: %s (ID: %d)", space.Name, space.ID)
	} else {
		// Get existing space
		space = &models.Space{}
		if err := models.DB.First(space).Error; err != nil {
			log.Printf("Failed to get existing space: %v", err)
			return
		}
		log.Printf("Using existing space: %s (ID: %d)", space.Name, space.ID)
	}

	// Check if we need to create default meal types
	var mealTypeCount int64
	models.DB.Model(&models.MealType{}).Where("space_id = ?", space.ID).Count(&mealTypeCount)
	log.Printf("Meal type count for space %d: %d", space.ID, mealTypeCount)

	if mealTypeCount == 0 {
		log.Println("Creating default meal types...")
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
				Space:     *space,
			}

			if err := models.DB.Create(&mealType).Error; err != nil {
				log.Printf("Failed to create meal type %s: %v", mt.name, err)
				continue
			}
			log.Printf("Created meal type: %s", mt.name)
		}
	}

	// Check if we need to create test data
	var foodCount int64
	models.DB.Model(&models.Food{}).Where("space_id = ?", space.ID).Count(&foodCount)
	log.Printf("Food count for space %d: %d", space.ID, foodCount)

	if foodCount == 0 {
		log.Println("Creating test foods...")
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

	// Check if we need to create recipe book
	var bookCount int64
	models.DB.Model(&models.RecipeBook{}).Where("space_id = ?", space.ID).Count(&bookCount)
	log.Printf("RecipeBook count for space %d: %d", space.ID, bookCount)

	if bookCount == 0 {
		log.Println("Creating test recipe book...")
		recipeBook, err := models.CreateRecipeBook(&user, space, "My Recipes", "A collection of my favorite recipes")
		if err != nil {
			log.Printf("Failed to create recipe book: %v", err)
		} else {
			log.Printf("Test recipe book created: %s (ID: %d)", recipeBook.Name, recipeBook.ID)
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
