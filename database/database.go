package database

import (
	"toncap-backend/config"
	"toncap-backend/logger"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	database, err := gorm.Open(sqlite.Open(config.DATABASE_FILE), &gorm.Config{})
	if err != nil {
		logger.Log.Fatalf("[database] connect error: %s", err.Error())
	}

	//	err = database.AutoMigrate(&types.Price{})
	//	if err != nil {
	//		return err
	//	}

	DB = database
	logger.Log.Info("[database] connected")
}
