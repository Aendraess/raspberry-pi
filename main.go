package main

import (
	_ "api/docs"
	"api/server"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Failed to load environment variables.")
	}
	server.InitalizeServer()
}
