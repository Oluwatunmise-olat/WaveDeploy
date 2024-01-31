package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
)

var DB *gorm.DB

func Connect() (*gorm.DB, error) {
	dsn := os.Getenv("MYSQL_CONNECTION_STRING")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(getQueryLogMode()),
	})
	if err != nil {
		return nil, err
	}
	DB = db

	return DB, nil
}

func getQueryLogMode() logger.LogLevel {
	if os.Getenv("ENVIRONMENT") == "production" {
		return logger.Silent
	}
	return logger.Info
}
