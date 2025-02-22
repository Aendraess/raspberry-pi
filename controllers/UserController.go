package controllers

import (
	"api/database"
	"api/models"
	"log"

	"github.com/gofiber/fiber/v2"
)

type UserController struct{}

func (uc *UserController) RegisterRoutes(app fiber.Router) {
	log.Println("Setting up user logs...")
	group := app.Group("/users")
	group.Post("/", uc.CreateUser)
	group.Get("/", uc.GetUsers)
}

// @Summary Get a list of users
// @Description Get a list of all users
// @Produce json
// @Tags User
// @Success 200 {array} models.User
// @Router /api/users [get]
func (uc *UserController) GetUsers(c *fiber.Ctx) error {
	var users []models.User

	database.DB.Find(&users)
	return c.JSON(users)
}

// @Summary Create a new user
// @Description Create a new user with the given name
// @Accept json
// @Produce json
// @Tags User
// @Param user body dtos.CreateUserRequest true "User object"
// @Success 200 {object} dtos.UserResponse
// @Router /api/users [post]
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

// @Summary Update a user
// @Description Update an existing user by ID
// @Accept json
// @Produce json
// @Tags User
// @Param id path int true "User ID"
// @Param user body dtos.UpdateUserRequest true "Updated user object"
// @Success 200 {object} dtos.UserResponse
// @Router /api/users/{id} [put]
func (uc *UserController) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User

	if err := database.DB.First(&user, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	var updateData models.User
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	database.DB.Model(&user).Updates(updateData)

	return c.JSON(user)
}
