package controllers

import (
	"api/database"
	"api/models"
	"log"

	"github.com/gofiber/fiber/v2"
)

type CategoryController struct{}

func (c *CategoryController) RegisterRoutes(app fiber.Router) {
	log.Println("Setting up category logs...")
	group := app.Group("/category")
	group.Get("/", c.GetCategories)
	group.Post("/", c.CreateCategory)
}

// @Summary Get a list of market categories
// @Description Get a list of all market categories
// @Produce json
// @Tags Category
// @Success 200 {array} models.Category
// @Router /api/category [get]
func (uc *CategoryController) GetCategories(c *fiber.Ctx) error {
	var marketItems []models.MarketItem

	database.DB.Find(&marketItems)
	return c.JSON(marketItems)
}

// @Summary Create a new market Category
// @Description Create a new market Category
// @Accept json
// @Produce json
// @Tags Category
// @Param user body dtos.CreateCategory true "MarketCategory object"
// @Success 200 {object} models.Category
// @Router /api/category [post]
func (uc *CategoryController) CreateCategory(c *fiber.Ctx) error {
	var marketItem models.MarketItem
	// Parse the request body into the User struct
	if err := c.BodyParser(&marketItem); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	database.DB.Create(&marketItem)

	return c.JSON(marketItem)
}