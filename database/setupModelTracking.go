package database

import (
	"api/models"

	"gorm.io/gorm"
)

func setupModelTracking() {
	// Register callbacks for all models
	DB.Callback().Create().After("gorm:create").Register("track_model_updates", func(db *gorm.DB) {
		if db.Error == nil {
			modelName := db.Statement.Schema.Name
			if modelName != "ModelUpdates" {
				updateModelTracking(db, modelName, "create")
			}
		}
	})

	DB.Callback().Update().After("gorm:update").Register("track_model_updates", func(db *gorm.DB) {
		if db.Error == nil {
			modelName := db.Statement.Schema.Name
			if modelName != "ModelUpdates" {
				updateModelTracking(db, modelName, "update")
			}
		}
	})	

	DB.Callback().Delete().After("gorm:delete").Register("track_model_updates", func(db *gorm.DB) {
		if db.Error == nil {
				modelName := db.Statement.Schema.Name
			if modelName != "ModelUpdates" {
				updateModelTracking(db, modelName, "delete")
			}
		}
	})
}

func updateModelTracking(db *gorm.DB, modelName, method string) {
	modelUpdate := models.ModelUpdates{
		ModelName: modelName,
		Method:    method,
	}
	
	// Upsert the model update record
	db.Where(models.ModelUpdates{ModelName: modelName}).
		Assign(models.ModelUpdates{Method: method}).
		FirstOrCreate(&modelUpdate)
}