package main

import (
	"toncap-backend/config"
	"toncap-backend/logger"
	"toncap-backend/router"
)

func main() {
	app := router.Router()

	logger.Log.Info("[http] started on: ", config.SERVER_ADDRESS)
	err := app.Run(config.SERVER_ADDRESS)
	if err != nil {
		logger.Log.Fatal("[http] error: ", err.Error())
	}
}
