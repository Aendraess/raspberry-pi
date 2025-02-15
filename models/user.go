package models

import (
	"time"
)

//BaseModel contains common fields for all models
type BaseModel struct {
	ID        int        `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt,omitempty" gorm:"index"`
}

type User struct {
	BaseModel
	Name string `json:"name"`
}
