package dtos

import "time"

type ChatMessageResponse struct {
	ID        int       `json:"id"`
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}

type ChatThreadResponse struct {
	ID        int                   `json:"id"`
	Title     string                `json:"title"`
	CreatedAt time.Time             `json:"createdAt"`
	UpdatedAt time.Time             `json:"updatedAt"`
	Messages  []ChatMessageResponse `json:"messages,omitempty"`
}

type AddMessageResponse struct {
	Reply string `json:"reply"`
}
