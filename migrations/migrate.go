package migrations

import (
	"github.com/Egorpalan/avito-shop/internal/models"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	// Автомиграции для моделей
	if err := db.AutoMigrate(&models.User{}, &models.Merch{}, &models.Transaction{}, &models.Inventory{}); err != nil {
		return err
	}

	// Заполнение таблицы merches начальными данными
	if err := SeedMerches(db); err != nil {
		return err
	}

	return nil
}
