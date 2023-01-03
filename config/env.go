package config

import (
	"os"
	"strconv"
	"time"
	"toncap-backend/logger"

	"github.com/joho/godotenv"
)

var (
	SERVER_ADDRESS       string
	CONTRACT_ADDRESS     string
	DATABASE_FILE        string
	ADS_FILE             string
	JETTONS_FILE         string
	JWT_SECRET           string
	GIN_MODE             string
	CURRENCY_UPDATE_TIME time.Duration
	ADS_UPDATE_TIME      time.Duration
	JETTONS_UPDATE_TIME  time.Duration
)

func init() {
	err := godotenv.Load()
	if err != nil {
		logger.Log.Fatalf("[env] .env load error: %s", err.Error())
	}
	logger.Log.Info("[env] loaded")

	SERVER_ADDRESS = os.Getenv("SERVER_ADDRESS")
	CONTRACT_ADDRESS = os.Getenv("CONTRACT_ADDRESS")
	DATABASE_FILE = os.Getenv("DATABASE_FILE")
	ADS_FILE = os.Getenv("ADS_FILE")
	JETTONS_FILE = os.Getenv("JETTONS_FILE")
	JWT_SECRET = os.Getenv("JWT_SECRET")
	GIN_MODE = os.Getenv("GIN_MODE")

	CURRENCY_UPDATE_TIME_SECS, _ := strconv.Atoi(os.Getenv("CURRENCY_UPDATE_TIME_SECS"))
	CURRENCY_UPDATE_TIME = time.Duration(CURRENCY_UPDATE_TIME_SECS) * time.Second

	ADS_UPDATE_TIME_SECS, _ := strconv.Atoi(os.Getenv("ADS_UPDATE_TIME_SECS"))
	ADS_UPDATE_TIME = time.Duration(ADS_UPDATE_TIME_SECS) * time.Second

	JETTONS_UPDATE_TIME_SECS, _ := strconv.Atoi(os.Getenv("JETTONS_UPDATE_TIME_SECS"))
	JETTONS_UPDATE_TIME = time.Duration(JETTONS_UPDATE_TIME_SECS) * time.Second
}
