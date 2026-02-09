package dtos

type CreateLogBookEntryRequest struct {
	Message  string `json:"message"`
	Level    string `json:"level"`
	Category string `json:"category"`
}
