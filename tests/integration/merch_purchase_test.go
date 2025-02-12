package integration

import (
	"log"
	"testing"

	"github.com/Egorpalan/avito-shop/internal/models"
	"github.com/Egorpalan/avito-shop/internal/repository"
	"github.com/Egorpalan/avito-shop/internal/service"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Подключаем тестовую БД
func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Println(err)
	}
	err = db.AutoMigrate(&models.User{}, &models.Transaction{}, &models.Merch{}, &models.Inventory{})
	if err != nil {
		log.Println(err)
	}
	return db
}

func TestMerchPurchaseScenario(t *testing.T) {
	db := setupTestDB()
	transactionRepo := repository.NewTransactionRepository(db)
	userRepo := repository.NewUserRepository(db)
	merchRepo := repository.NewMerchRepository(db)
	merchService := service.NewMerchService(merchRepo, transactionRepo, userRepo) // Используем MerchService

	user := models.User{Username: "testuser", Coins: 500}
	db.Create(&user)

	merch := models.Merch{Name: "T-Shirt", Price: 300}
	db.Create(&merch)

	err := merchService.BuyMerchByUsername("testuser", merch.ID)
	assert.NoError(t, err)

	var updatedUser models.User
	db.First(&updatedUser, "username = ?", "testuser")
	assert.Equal(t, 200, updatedUser.Coins)

	var inventory models.Inventory
	db.First(&inventory, "user_id = ? AND merch_id = ?", updatedUser.ID, merch.ID)
	assert.Equal(t, 1, inventory.Quantity)

	t.Log("TestMerchPurchaseScenario complete")
}
