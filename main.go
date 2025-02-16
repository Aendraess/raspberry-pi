package main

import (
	_ "api/docs"
	"api/server"
)

// @title User API
// @version 1.0
// @description This is a sample API for managing users.
// @host localhost:8080
// @BasePath /
func main() {
	// Initialize the database
	server.InitalizeServer()
}
