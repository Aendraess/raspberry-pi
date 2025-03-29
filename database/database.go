package database

import (
	"log"

	"gorm.io/driver/sqlite" // SQLite driver for GORM
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	// Connect to SQLite database (you can replace this with any other database driver)
	DB, err = gorm.Open(sqlite.Open("./database.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	DB = DB.Set("gorm:query_option", "WHERE deleted_at IS NULL")
	
	
	
	// Auto migrate 
	migrateDb()
	
	// Set up callbacks for tracking model updates
	setupModelTracking()

	log.Println("Database connected and migrated successfully")
}
