package controllers

import (
	"api/database"
	"api/models"
	"api/server"

	"github.com/gofiber/fiber/v2"
)

type UserController struct{}

func (uc *UserController) RegisterRoutes() {
	group := server.App.Group("/users")
	group.Post("/", uc.CreateUser)
	group.Get("/", uc.GetUsers)
}

// @Summary Get a list of users
// @Description Get a list of all users
// @Produce json
// @Success 200 {array} models.User
// @Router /users [get]
func (uc *UserController) GetUsers(c *fiber.Ctx) error {
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
func (uc *UserController) CreateUser(c *fiber.Ctx) error {
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
