package server

import (
	"api/controllers"
	"api/database"
	"os"
	"log"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	"github.com/swaggo/swag"
)

var App *fiber.App
var Port = ""
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Title:            "User API",
	Description:      "This is a sample API for managing users.",
	Host:            "",
	BasePath:         "/",
	Schemes:          []string{"http"},
	InfoInstanceName: "swagger",
}

// InitializeSwagger sets up Swagger with a dynamic host.
func InitializeSwagger() {
	port := os.Getenv("GOPORT")
	log.Println("PORT: ",port)
	if port == "" {
		port = "8080"
	}
	Port = port
	SwaggerInfo.Host = fmt.Sprintf("localhost:%s", port)
	swag.Register(SwaggerInfo.InfoInstanceName, SwaggerInfo)
}
func InitalizeServer() {
	// Create a new Fiber app
	App = fiber.New()
	InitializeSwagger()
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

	App.Listen(":"+Port)
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
