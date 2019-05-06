package main

import (
	"fmt"
	"go-booking/configure"
	"go-booking/controllers"
	"go-booking/model/Orders"
	"go-booking/routing"
)

func main() {
	server := configure.ServerConfig{}
	err := server.LoadConfigServer()
	if err != nil {
		fmt.Println("CHECK CONFIG.JSON AND ENVIRONMENT")
	}

	controllers.GetConfigJWT()
	InitApp()

	router := routing.ConfigureRouters()
	router.Run(":" + server.Port)
}

func InitApp() {
	Orders.Token366.GetToken()
}
