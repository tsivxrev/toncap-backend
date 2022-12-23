package main

import (
	"log"
	"time"
	"toncap-backend/db"
	"toncap-backend/router"
	"toncap-backend/storage"
	"toncap-backend/tasks"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	err = storage.Init()
	if err != nil {
		log.Fatal("Storage init error: " + err.Error())
	}
}

func main() {
	fetchAdsTask := func(stop chan bool) {
		ads, err := db.GetAds()
		if err != nil {
			log.Println("Ads fetch error: " + err.Error())
		}

		storage.App.Ads = ads
	}

	fetchExchangeRateTask := func(stop chan bool) {
		exchangeRate, err := db.GetExchangeRate()
		if err != nil {
			log.Println("Exchange rate fetch error: " + err.Error())
		}

		storage.App.ExchangeRate = exchangeRate
	}

	// run once and add to tasks
	fetchAdsTask(nil)
	fetchExchangeRateTask(nil)

	tasks.New(fetchAdsTask, 1*time.Minute)
	tasks.New(fetchExchangeRateTask, 1*time.Minute)

	app := router.Router()
	app.Run(storage.App.ServerAddress)
}
