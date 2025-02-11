package db

import (
	"github.com/Egorpalan/avito-shop/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbConn *gorm.DB

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := "host=" + cfg.DBHost + " user=" + cfg.DBUser + " password=" + cfg.DBPassword + " dbname=" + cfg.DBName + " port=" + cfg.DBPort + " sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	dbConn = db
	return db, nil
}

// GetDB возвращает подключение к базе данных
func GetDB() *gorm.DB {
	return dbConn
}
