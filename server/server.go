package server

import (
	"api/controllers"
	"api/database"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// @title Andreas API
// @version 1.0
// @description Andreas Personal API for home stuff and testing.
// @termsOfService http://swagger.io/terms/

// @contact.name Andreas Löfkvist
// @contact.url http://localhost
// @contact.email andreasmlofkvist@gmail.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8081
// @BasePath /api
// @schemes http
var App *fiber.App
var Api fiber.Router
var Port = ""

func InitalizeServer() {
	// Create a new Fiber app
	Port = os.Getenv("GOPORT")
	App = fiber.New()
	// Add CORS middleware
	App.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Allows all origins
	}))
	database.InitDB()
	Api = App.Group("/api")
	SetupRoutes(&Api)
	// Serve Swagger UI
	App.Get("/swagger/*", fiberSwagger.WrapHandler)
	log.Println("Registered Routes:")
	for _, route := range App.Stack() {
		for _, r := range route {
			log.Printf("%s %s", r.Method, r.Path)
		}
	}

	App.Listen(":" + Port)
}

// SetupRoutes automatically registers controllers
func SetupRoutes(app *fiber.Router) {
	controllersList := []controllers.Controller{
		&controllers.UserController{},
		&controllers.CategoryController{},
		&controllers.MarketItemController{},
	}

	for _, controller := range controllersList {
		controller.RegisterRoutes(*app)
	}
}
