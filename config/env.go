package config

import (
	"encoding/json"
	"os"
	"strconv"
	"time"
	"toncap-backend/logger"
	"toncap-backend/types"

	"github.com/joho/godotenv"
)

var (
	SERVER_ADDRESS       string
	CONTRACT_ADDRESS     string
	DATABASE_FILE        string
	ADS_FILE             string
	JETTONS_FILE         string
	HASH_SECRET          string
	GIN_MODE             string
	CURRENCY_UPDATE_TIME time.Duration
	ADS_UPDATE_TIME      time.Duration
	JETTONS_UPDATE_TIME  time.Duration

	JETTONS map[string]types.Jetton
	ADS     []types.Ad
)

func FetchAds() {
	adsFile, err := os.ReadFile(ADS_FILE)
	if err != nil {
		logger.Log.Errorf("Cannot read ads file: %s", err.Error())
		return
	}

	var ads []types.Ad
	err = json.Unmarshal(adsFile, &ads)
	if err != nil {
		logger.Log.Errorf("Cannot read ads file: %s", err.Error())
		return
	}

	logger.Log.Info("Ads fetched")
	ADS = ads
}

func FetchJettons() {
	jettonsFile, err := os.ReadFile(JETTONS_FILE)
	if err != nil {
		logger.Log.Errorf("Cannot read jettons file: %s", err.Error())
		return
	}

	var jettons map[string]types.Jetton
	err = json.Unmarshal(jettonsFile, &jettons)
	if err != nil {
		logger.Log.Errorf("Cannot parse jettons file: %s", err.Error())
		return
	}

	logger.Log.Info("Jettons fetched")
	JETTONS = jettons
}

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
	HASH_SECRET = os.Getenv("HASH_SECRET")
	GIN_MODE = os.Getenv("GIN_MODE")

	CURRENCY_UPDATE_TIME_SECS, _ := strconv.Atoi(os.Getenv("CURRENCY_UPDATE_TIME_SECS"))
	CURRENCY_UPDATE_TIME = time.Duration(CURRENCY_UPDATE_TIME_SECS) * time.Second

	FetchAds()
	FetchJettons()

	//ADS_UPDATE_TIME_SECS, _ := strconv.Atoi(os.Getenv("ADS_UPDATE_TIME_SECS"))
	//ADS_UPDATE_TIME = time.Duration(ADS_UPDATE_TIME_SECS) * time.Second

	//JETTONS_UPDATE_TIME_SECS, _ := strconv.Atoi(os.Getenv("JETTONS_UPDATE_TIME_SECS"))
	//JETTONS_UPDATE_TIME = time.Duration(JETTONS_UPDATE_TIME_SECS) * time.Second
}
