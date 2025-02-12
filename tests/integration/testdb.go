package integration

import (
	"log"

	"github.com/Egorpalan/avito-shop/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	err = db.AutoMigrate(&models.User{}, &models.Transaction{}, &models.Merch{}, &models.Inventory{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	seedMerches(db)

	return db
}

func seedMerches(db *gorm.DB) {
	merches := []models.Merch{
		{Name: "t-shirt", Price: 80},
		{Name: "cup", Price: 20},
		{Name: "book", Price: 50},
		{Name: "pen", Price: 10},
		{Name: "powerbank", Price: 200},
		{Name: "hoody", Price: 300},
		{Name: "umbrella", Price: 200},
		{Name: "socks", Price: 10},
		{Name: "wallet", Price: 50},
		{Name: "pink-hoody", Price: 500},
	}

	for _, merch := range merches {
		db.FirstOrCreate(&merch, models.Merch{Name: merch.Name})
	}
}

func clearTestDB(db *gorm.DB) {
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM transactions")
	db.Exec("DELETE FROM merches")
	db.Exec("DELETE FROM inventories")
}
