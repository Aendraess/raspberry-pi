package models

// ModelUpdates tracks the last update method for each model in the system
type ModelUpdates struct {
	BaseModel
	ModelName string `json:"model_name" gorm:"uniqueIndex"`
	Method    string `json:"action"`
}
