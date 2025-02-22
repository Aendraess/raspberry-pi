package models

type MarketItem struct {
	BaseModel
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	Title       string  `json:"title"`
	CategoryId  uint
	Category    Category `json:"category" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}