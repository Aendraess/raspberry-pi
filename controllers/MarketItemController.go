package controllers

import (
	"api/database"
	"api/models"
	"log"

	"github.com/gofiber/fiber/v2"
)

type MarketItemController struct{}

func (mc *MarketItemController) RegisterRoutes(app fiber.Router) {
	log.Println("Setting up user logs...")
	group := app.Group("/marketItem")
	group.Post("/", mc.CreateMarketItem)
	group.Get("/", mc.GetMarketItems)
	group.Delete("/", mc.DeleteMarketItem)
	group.Put("/", mc.UpdateMarketItem)
	group.Get("/category", mc.GetMarketCategories)
	group.Post("/category", mc.CreateMarketCategory)

}

// @Summary Get a list of market items
// @Description Get a list of all market items
// @Produce json
// @Tags MarketItem
// @Success 200 {array} models.MarketItem
// @Router /api/marketItem [get]
func (uc *MarketItemController) GetMarketItems(c *fiber.Ctx) error {
	var marketItems []models.MarketItem

	database.DB.Find(&marketItems)
	return c.JSON(marketItems)
}


// @Summary Get a list of market categories
// @Description Get a list of all market categories
// @Produce json
// @Tags MarketItem
// @Success 200 {array} models.MarketCategory
// @Router /api/marketItem/category [get]
func (uc *MarketItemController) GetMarketCategories(c *fiber.Ctx) error {
	var marketItems []models.MarketItem

	database.DB.Find(&marketItems)
	return c.JSON(marketItems)
}

// @Summary Create a new market item
// @Description Create a new market item
// @Accept json
// @Produce json
// @Tags MarketItem
// @Param user body dtos.CreateMarketItemRequest true "MarketItem object"
// @Success 200 {object} models.MarketItem
// @Router /api/users [post]
func (uc *MarketItemController) CreateMarketItem(c *fiber.Ctx) error {
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

// @Summary Create a new market Category
// @Description Create a new market Category
// @Accept json
// @Produce json
// @Tags MarketItem
// @Param user body dtos.CreateCategory true "MarketCategory object"
// @Success 200 {object} models.MarketCategory
// @Router /api/marketItem/Category [post]
func (uc *MarketItemController) CreateMarketCategory(c *fiber.Ctx) error {
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

// @Summary Delete a market item
// @Description Delete a market item
// @Accept json
// @Produce json
// @Tags MarketItem
// @Param id path int true "Market item ID"
// @Success 200 {object} models.MarketItem
// @Router /api/marketItem [delete]
func (uc *MarketItemController) DeleteMarketItem(c *fiber.Ctx) error {
	var marketItem models.MarketItem
	id := c.Params("id")
	if err := database.DB.First(&marketItem, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Market item not found",
		})
	}
	err := database.DB.Delete(marketItem)
	if (err != nil) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Failed to delete item",
		})
	}
	return c.JSON(marketItem)
	}

// @Summary Update a market item
// @Description Update an existing market item
// @Accept json
// @Produce json
// @Tags MarketItem
// @Param id path int true "MarketItemId ID"
// @Param user body dtos.UpdateMarketItemRequest true "Updated market item object"
// @Success 200 {object} models.MarketItem
// @Failure 400 {object} fiber.Map "Bad Request"
// @Failure 404 {object} fiber.Map "Market Item Not Found"
// @Router /api/marketItem/{id} [put]

func (uc *MarketItemController) UpdateMarketItem(c *fiber.Ctx) error {
	id := c.Params("id")
	var marketItem models.MarketItem

	if err := database.DB.First(&marketItem, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Market item not found",
		})
	}

	var updateData models.MarketItem
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	database.DB.Model(&marketItem).Updates(updateData)

	return c.JSON(marketItem)
}
