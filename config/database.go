package config

import (
	"godoor/models"
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	var err error

	// Get database path from environment variable or use default
	dbPath := os.Getenv("GODOOR_DB_PATH")
	if dbPath == "" {
		dbPath = "godoor.db"
	}

	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Set global DB variable
	models.DB = DB
}
