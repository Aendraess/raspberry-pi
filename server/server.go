package server

import (
	"api/controllers"
	"api/database"
	"api/mcpServer"
	"api/services"
	"log"
	"os"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	mcp "github.com/mark3labs/mcp-go/server"
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
	mcpSrv := mcpServer.NewServer()
	mcpHTTP := mcp.NewStreamableHTTPServer(mcpSrv, mcp.WithEndpointPath("/mcp"))
	App.All("/mcp/*", adaptor.HTTPHandler(mcpHTTP))
	Api = App.Group("/api")
	chatService := services.NewChatService(mcpSrv)
	SetupRoutes(&Api, chatService)
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
func SetupRoutes(app *fiber.Router, chatService *services.ChatService) {
	controllersList := []controllers.Controller{
		&controllers.UserController{},
		&controllers.CategoryController{},
		&controllers.MarketItemController{},
		&controllers.ApiKeyController{},
		&controllers.BloodPressureController{},
		controllers.NewChatController(chatService),
	}

	for _, controller := range controllersList {
		controller.RegisterRoutes(*app)
	}
}
