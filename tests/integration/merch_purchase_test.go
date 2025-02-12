package integration

import (
	"testing"

	"github.com/Egorpalan/avito-shop/internal/models"
	"github.com/Egorpalan/avito-shop/internal/repository"
	"github.com/Egorpalan/avito-shop/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestMerchPurchaseScenario(t *testing.T) {
	db := setupTestDB()
	defer clearTestDB(db)

	userRepo := repository.NewUserRepository(db)
	merchRepo := repository.NewMerchRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)

	userService := service.NewUserService(userRepo)
	merchService := service.NewMerchService(merchRepo, transactionRepo, userRepo)

	user := &models.User{Username: "Alice", Password: "password123"}
	err := userService.Register(user)
	assert.NoError(t, err, "User registration should succeed")

	user, err = userRepo.GetUserByUsername("Alice")
	assert.NoError(t, err)

	savedMerch, err := merchRepo.GetMerchByName("t-shirt")
	assert.NoError(t, err)
	assert.NotNil(t, savedMerch, "Merch should be found in DB")

	err = merchService.BuyMerchByUsername("Alice", savedMerch.ID)
	assert.NoError(t, err, "Purchase should succeed")

	newBalance, err := transactionRepo.GetUserBalance(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, 920, newBalance, "User's balance should be 920 after purchase")

	var inventory models.Inventory
	err = db.Where("user_id = ? AND merch_id = ?", user.ID, savedMerch.ID).First(&inventory).Error
	assert.NoError(t, err, "Inventory should contain the purchased merch")
	assert.Equal(t, 1, inventory.Quantity, "User should have 1 t-shirt in inventory")

	t.Log("TestMerchPurchaseScenario complete")
}
