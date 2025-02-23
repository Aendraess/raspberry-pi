package database

import (
	"api/models"
	"log"
)

func migrateDb() {
	err := DB.AutoMigrate(
		&models.User{}, &models.ApiKey{}, &models.MarketItem{}, &models.Category{})
	if err != nil {
		log.Fatal("Failed to migrate user, ", err)
	}
}
