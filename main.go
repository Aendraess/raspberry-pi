package main

import (
	_ "api/docs"
	"api/server"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Initialize the database
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "development"
	}

	// Load the corresponding .env file
	envFile := ".env." + env
	err := godotenv.Load(envFile)
	if err != nil {
		log.Println("Failed to load environment variables.")
	}
	server.InitalizeServer()
}
