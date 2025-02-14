package migrations

import (
	"github.com/Egorpalan/avito-shop/internal/models"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&models.User{}, &models.Merch{}, &models.Transaction{}, &models.Inventory{}); err != nil {
		return err
	}

	if err := SeedMerches(db); err != nil {
		return err
	}

	if err := AddIndexes(db); err != nil {
		return err
	}

	return nil
}
