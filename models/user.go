package models

type User struct {
	BaseModel
	Name  string `json:"name"`
	Phone string `json:"phoneNumber"`
	Email string `json:"email"`
}
