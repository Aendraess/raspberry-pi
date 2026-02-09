package models

import "time"

type LogBookEntry struct {
	BaseModel
	Message   string    `json:"message"`
	Level     string    `json:"level"`
	Category  string    `json:"category"`
	Timestamp time.Time `json:"timestamp"`
}
