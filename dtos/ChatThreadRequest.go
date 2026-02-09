package dtos

type CreateThreadRequest struct {
	Title string `json:"title"`
}

type AddMessageRequest struct {
	Message string `json:"message"`
}
