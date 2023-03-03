package database

import (
	"toncap-backend/types"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() error {
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&types.Ad{}, &types.Price{})
	if err != nil {
		return err
	}

	DB = db
	return nil
}
