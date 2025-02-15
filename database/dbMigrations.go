package database

import (
	"api/models"
	"log"
)

func migrateDb() {
	err := DB.AutoMigrate(
		&models.User{})
	if err != nil {
		log.Fatal("Failed to migrate user, ", err)
	}
}
