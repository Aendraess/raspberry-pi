package models

type ChatMessage struct {
	BaseModel
	ThreadID int    `json:"threadId" gorm:"not null;index"`
	Role     string `json:"role" gorm:"size:32;not null"` // "user" | "assistant" | "system"
	Content  string `json:"content" gorm:"type:text;not null"`
}
