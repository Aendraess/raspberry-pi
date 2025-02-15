
package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "api/docs" // Import the generated docs
	fiberSwagger "github.com/swaggo/fiber-swagger"
	"gorm.io/driver/sqlite" // SQLite driver for GORM
	"gorm.io/gorm"
	"log"
	"api/models"
)

// Database connection
var db *gorm.DB

// @title User API
// @version 1.0
// @description This is a sample API for managing users.
// @host localhost:8080
// @BasePath /
func main() {
	// Initialize the database
	initDB()

	// Create a new Fiber app
	app := fiber.New()

	// Add CORS middleware
	app.Use(cors.New())

	// Serve Swagger UI
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// Define routes
	app.Get("/users", getUsers)
	app.Post("/users", createUser)

	// Start the server
	log.Println("Server running and listening on port: 8080")
	app.Listen(":8080")
}

// Initialize the database
func initDB() {
	var err error
	// Connect to SQLite database (you can replace this with any other database driver)
	db, err = gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate the User model
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Failed to auto migrate database:", err)
	}

	log.Println("Database connected and migrated successfully")
}

// @Summary Get a list of users
// @Description Get a list of all users
// @Produce json
// @Success 200 {array} models.User
// @Router /users [get]
func getUsers(c *fiber.Ctx) error {
	var users []models.User
	// Fetch all users from the database
	db.Find(&users)
	return c.JSON(users)
}

// @Summary Create a new user
// @Description Create a new user with the given name
// @Accept json
// @Produce json
// @Param user body models.User true "User object"
// @Success 200 {object} models.User
// @Router /users [post]
func createUser(c *fiber.Ctx) error {
	var user models.User
	// Parse the request body into the User struct
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// Create the user in the database
	db.Create(&user)

	return c.JSON(user)
}

