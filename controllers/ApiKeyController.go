package controllers

import (
	"api/database"
	"api/models"
	"log"

	"github.com/gofiber/fiber/v2"
)

type ApiKeyController struct{}

func (uc *ApiKeyController) RegisterRoutes(app fiber.Router) {
	log.Println("Setting up user logs...")
	group := app.Group("/api_keys")
	group.Post("/", uc.CreateApiKey)
	group.Get("/", uc.GetApiKeys)
}

// @Summary Get a list of ApiKeys
// @Description Get a list of all ApiKeys
// @Produce json
// @Tags ApiKey
// @Success 200 {array} models.ApiKey
// @Router /api/api_key [get]
func (uc *ApiKeyController) GetApiKeys(c *fiber.Ctx) error {
	var apiKeys []models.ApiKey

	database.DB.Find(&apiKeys)
	return c.JSON(apiKeys)
}

// @Summary Create a new API KEY
// @Description Create a new API key
// @Accept json
// @Produce json
// @Tags ApiKey
// @Param api_key body dtos.CreateApiKeyRequest true "ApiKey object"
// @Success 200 {object} models.ApiKey
// @Router /api/api_key [post]
func (uc *ApiKeyController) CreateApiKey(c *fiber.Ctx) error {
	var apiKey models.ApiKey
	// Parse the request body into the User struct
	if err := c.BodyParser(&apiKey); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	database.DB.Create(&apiKey)

	return c.JSON(apiKey)
}

// @Summary Update an API key
// @Description Update an existing APi key by ID
// @Accept json
// @Produce json
// @Tags ApiKey
// @Param id path int true "API Key ID"
// @Param user body dtos.CreateApiKeyRequest true "Updated user object"
// @Success 200 {object} models.ApiKey
// @Router /api/api_key/{id} [put]
func (uc *ApiKeyController) UpdateApiKey(c *fiber.Ctx) error {
	id := c.Params("id")
	var apiKey models.ApiKey

	if err := database.DB.First(&apiKey, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Api key not found",
		})
	}

	var updateData models.User
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	database.DB.Model(&apiKey).Updates(updateData)

	return c.JSON(apiKey)
}
