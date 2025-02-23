package dtos

type CreateApiKeyRequest struct {
	ClientSecret string `json:"client_secret"`
	ClientId     string `json:"client_id"`
	ApiKey       string `json:"api_key"`
	Type         string `json:"api_type"`
	UserId       uint   `json:"user_id"`
}