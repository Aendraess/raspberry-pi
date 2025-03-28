package dtos

type CreateMarketItemRequest struct {
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	Title       string  `json:"title"`
	CategoryId  uint    `json:"categoryId"`
	UserId      uint    `json:"userId"`
}

type UpdateMarketItemRequest struct {
	Id          string  `json:"id"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	Title       string  `json:"title"`
	CategoryId  uint    `json:"categoryId"`
	UserId      uint    `json:"userId"`
}

type CreateCategory struct {
	Title string `json:"title"`
}

type UpdateCategory struct {
	Title string `json:"title"`
}