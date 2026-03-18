package main

import (
	"haveYouWorkedOutToday/config"
	"haveYouWorkedOutToday/router"
)

func main() {
	config.InitConfig()

	r := router.SetupRouter()

	port := config.AppConfig.App.Port
	if port == "" {
		port = "8080"
	}

	r.Run(port) // listen and serve on 0.0.0.0:8080
}
