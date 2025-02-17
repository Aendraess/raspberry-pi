package dtos

import "time"

type CreateUserRequest struct {
	Name  string `json:"name"`
	Phone string `json:"phoneNumber"`
	Email string `json:"email"`
}

type UpdateUserRequest struct {
	id    int    `json:id`
	Name  string `json:"name"`
	Phone string `json:"phoneNumber"`
	Email string `json:"email"`
}

type UserResponse struct {
	Name      string    `json:"name"`
	Phone     string    `json:"phoneNumber"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}
