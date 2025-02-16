package server

import (
	"api/controllers"
	"api/database"
	"log"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

var App *fiber.App

func InitalizeServer() {
	// Create a new Fiber app
	App = fiber.New()

	// Add CORS middleware
	App.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Allows all origins
	}))
	database.InitDB()
	SetupRoutes(App)
	// Serve Swagger UI
	App.Get("/swagger/*", fiberSwagger.WrapHandler)
	log.Println("Registered Routes:")
	for _, route := range App.Stack() {
		for _, r := range route {
			log.Printf("%s %s", r.Method, r.Path)
		}
	}

	App.Listen(":8081")
}

// SetupRoutes automatically registers controllers
func SetupRoutes(app *fiber.App) {
	controllersList := []controllers.Controller{
		&controllers.UserController{},
	}

	for _, controller := range controllersList {
		controller.RegisterRoutes(app)
	}
}
