package server

import (
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
    AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
    AllowHeaders: "Origin, Content-Type, Accept, Authorization",
}))


	// Serve Swagger UI
	App.Get("/swagger/*", fiberSwagger.WrapHandler)

	// Start the server
	log.Println("Server running and listening on port: 8081")
	App.Listen(":8081")
}
