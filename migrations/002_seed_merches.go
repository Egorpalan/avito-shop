package migrations

import (
	"github.com/Egorpalan/avito-shop/internal/models"
	"gorm.io/gorm"
)

func SeedMerches(db *gorm.DB) error {
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
		if err := db.FirstOrCreate(&merch, models.Merch{Name: merch.Name}).Error; err != nil {
			return err
		}
	}

	return nil
}
