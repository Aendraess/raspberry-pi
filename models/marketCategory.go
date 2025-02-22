package models

type Category struct {
	BaseModel
	Title string `json:"title"`
}