package main

import (
	_ "api/docs"
	"api/server"
)


func main() {
	// Initialize the database
	server.InitalizeServer()
}
