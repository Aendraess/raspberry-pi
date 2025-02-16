package routes

import (
	"api/controllers"
)

// SetupRoutes automatically registers controllers
func SetupRoutes() {
	controllersList := []controllers.Controller{
		&controllers.UserController{},
	}

	for _, controller := range controllersList {
		controller.RegisterRoutes()
	}
}
