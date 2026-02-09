package models

type ChatThread struct {
	BaseModel
	Title string `json:"title" gorm:"size:512"`
}
