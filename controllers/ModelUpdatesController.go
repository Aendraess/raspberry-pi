package controllers

import (
	"api/database"
	"api/models"
	"log"

	"github.com/gofiber/fiber/v2"
)

type ModelUpdatesController struct{}

func (mc *ModelUpdatesController) RegisterRoutes(app fiber.Router) {
	log.Println("Setting up user logs...")
	group := app.Group("/marketitem")
	group.Get("/", mc.GetModelUpdates)
}

// @Summary Get Model updates
// @Description Get a list of all market items
// @Produce json
// @Tags ModelUpdates
// @Success 200 {array} models.ModelUpdates
// @Router /api/modelUpdates [get]
func (uc *ModelUpdatesController) GetModelUpdates(c *fiber.Ctx) error {
	var modelUpdates []models.ModelUpdates

	database.DB.Find(&modelUpdates)
	return c.JSON(modelUpdates)
}