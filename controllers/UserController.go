package controllers

import (
	"api/database"
	"api/models"
	"api/server"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes() {
	group := server.App.Group("/users")
	group.Post("/", CreateUser)
	group.Get("/", GetUsers)
}

// @Summary Get a list of users
// @Description Get a list of all users
// @Produce json
// @Success 200 {array} models.User
// @Router /users [get]
func GetUsers(c *fiber.Ctx) error {
	var users []models.User

	database.DB.Find(&users)
	return c.JSON(users)
}

// @Summary Create a new user
// @Description Create a new user with the given name
// @Accept json
// @Produce json
// @Param user body models.User true "User object"
// @Success 200 {object} models.User
// @Router /users [post]
func CreateUser(c *fiber.Ctx) error {
	var user models.User
	// Parse the request body into the User struct
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	database.DB.Create(&user)

	return c.JSON(user)
}
