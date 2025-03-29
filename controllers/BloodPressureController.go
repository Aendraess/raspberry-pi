package controllers

import (
	"api/database"
	"api/dtos"
	"api/models"
	"log"

	"github.com/gofiber/fiber/v2"
)

type BloodPressureController struct {}

func (uc *BloodPressureController) RegisterRoutes(app fiber.Router) {
	log.Println("Setting up user logs...")
	group := app.Group("/blood_pressure")
	group.Post("/", uc.CreateBloodPressure)
	group.Get("/", uc.GetBloodPressureReports)
}

// @Summary Get a list of blood pressures
// @Description Get a list of all blood pressure recordings
// @Produce json
// @Tags BloodPressure
// @Success 200 {array} models.BloodPressure
// @Router /api/blood_pressure [get]
func (uc *BloodPressureController) GetBloodPressureReports(c *fiber.Ctx) error {
	var BloodPressureReports []models.BloodPressure

	database.DB.Find(&BloodPressureReports)
	return c.JSON(BloodPressureReports)
}

// @Summary Create a new Blood Pressure report
// @Description Create a new Blood Pressure report
// @Accept json
// @Produce json
// @Tags BloodPressure
// @Param user body dtos.CreateBloodPressure true "BloodPressure object"
// @Success 200 {object} models.BloodPressure
// @Router /api/blood_pressure [post]
func (uc *BloodPressureController) CreateBloodPressure(c *fiber.Ctx) error {
	var BloodPressureReport dtos.CreateBloodPressure
	// Parse the request body into the User struct
	if err := c.BodyParser(&BloodPressureReport); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}
	var BloodPressure = models.BloodPressure{
		Systolic: BloodPressureReport.Systolic,
		Diastolic: BloodPressureReport.Diastolic,
		Pulse: BloodPressureReport.Pulse,
		Medicine: BloodPressureReport.Medicine,
	}
	database.DB.Create(&BloodPressure)

	return c.JSON(BloodPressure)
}