package storage

import (
	"os"
	"toncap-backend/types"
)

var App *types.AppStorage

func Init() error {
	App = &types.AppStorage{
		ServerAddress:   os.Getenv("SERVER_ADDRESS"),
		ContractAddress: os.Getenv("CONTRACT_ADDRESS"),
		AdsFilePath:     os.Getenv("ADS_FILE_PATH"),
	}

	return nil
}
