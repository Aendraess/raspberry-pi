package database

import (
	"api/models"
	"log"
)

func migrateDb() {
	err := DB.AutoMigrate(
		&models.User{}, &models.ApiKey{}, &models.MarketItem{}, &models.Category{}, &models.BloodPressure{}, &models.ModelUpdates{})
	if err != nil {
		log.Fatal("Failed to migrate, ", err)
	}
}
