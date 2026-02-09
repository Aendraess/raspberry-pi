package controllers

import (
	"api/database"
	"api/models"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type LogBookEntryController struct{}

func (uc *LogBookEntryController) RegisterRoutes(app fiber.Router) {

	group := app.Group("/log_book")
	group.Post("/", uc.CreateLogBookEntry)
	group.Get("/", uc.GetLogBookEntries)
	group.Delete("/:id", uc.DeleteLogBookEntry)
}

// @Summary Get a list of LogBookEntries
// @Description Get a list of all LogBookEntries
// @Produce json
// @Tags LogBookEntry
// @Success 200 {array} models.LogBookEntry
// @Router /api/log_book [get]
func (uc *LogBookEntryController) GetLogBookEntries(c *fiber.Ctx) error {
	var logBookEntries []models.LogBookEntry

	database.DB.Find(&logBookEntries)
	return c.JSON(logBookEntries)
}

// @Summary Create a new LogBookEntry
// @Description Create a new LogBookEntry
// @Accept json
// @Produce json
// @Tags LogBookEntry
// @Param log_book body dtos.CreateLogBookEntryRequest true "LogBookEntry object"
// @Success 200 {object} models.LogBookEntry
// @Router /api/log_book [post]
func (uc *LogBookEntryController) CreateLogBookEntry(c *fiber.Ctx) error {
	var logBookEntry models.LogBookEntry
	// Parse the request body into the User struct
	if err := c.BodyParser(&logBookEntry); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse LogBookEntry JSON",
		})
	}
	logBookEntry.Timestamp = time.Now()

	database.DB.Create(&logBookEntry)

	return c.JSON(logBookEntry)
}

// @Summary Delete a LogBookEntry
// @Description Delete a LogBookEntry by ID
// @Accept json
// @Produce json
// @Tags LogBookEntry
// @Param id path int true "LogBookEntry ID"
// @Success 200 {object} models.LogBookEntry
// @Router /api/log_book/{id} [delete]
func (uc *LogBookEntryController) DeleteLogBookEntry(c *fiber.Ctx) error {
	id := c.Params("id")
	// Cast id to int
	idInt, err := strconv.Atoi(id)
	log.Println("Deleting LogBookEntry with ID: ", idInt)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}
	var logBookEntry models.LogBookEntry

	if err := database.DB.First(&logBookEntry, idInt).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "LogBookEntry not found",
		})
	}
	log.Println("LogBookEntry found: ", logBookEntry)

	database.DB.Delete(&logBookEntry)

	log.Println("LogBookEntry deleted.")

	return c.JSON(logBookEntry)
}
