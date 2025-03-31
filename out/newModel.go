package models

type NewModel struct {
	BaseModel
    name        string    `json:"name"`
    price       float32   `json:"price"`
    quantity    uint      `json:"quantity"`
    is_active   bool      `json:"is_active"`
    user        User      `json:"user,  gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;""`
    userId      uint      `json:"userId"`
}