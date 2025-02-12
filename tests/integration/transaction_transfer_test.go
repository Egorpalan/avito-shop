package integration

import (
	"testing"

	"github.com/Egorpalan/avito-shop/internal/models"
	"github.com/Egorpalan/avito-shop/internal/repository"
	"github.com/Egorpalan/avito-shop/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestTransferCoinsScenario(t *testing.T) {
	db := setupTestDB()
	defer clearTestDB(db)

	userRepo := repository.NewUserRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	userService := service.NewUserService(userRepo)
	transactionService := service.NewTransactionService(transactionRepo, userRepo)

	user1 := &models.User{Username: "Alice", Password: "password123"}
	user2 := &models.User{Username: "Bob", Password: "password123"}

	err := userService.Register(user1)
	assert.NoError(t, err, "User Alice should be registered")

	err = userService.Register(user2)
	assert.NoError(t, err, "User Bob should be registered")

	user1, err = userRepo.GetUserByUsername("Alice")
	assert.NoError(t, err)
	user2, err = userRepo.GetUserByUsername("Bob")
	assert.NoError(t, err)

	err = transactionService.SendCoins(user1.ID, user2.ID, 200)
	assert.NoError(t, err, "Transfer should succeed")

	newBalanceAlice, err := transactionRepo.GetUserBalance(user1.ID)
	assert.NoError(t, err)
	assert.Equal(t, 800, newBalanceAlice, "Alice's balance should be 800")

	newBalanceBob, err := transactionRepo.GetUserBalance(user2.ID)
	assert.NoError(t, err)
	assert.Equal(t, 1200, newBalanceBob, "Bob's balance should be 1200")

	var transaction models.Transaction
	db.First(&transaction, "from_user_id = ? AND to_user_id = ?", user1.ID, user2.ID)
	assert.Equal(t, 200, transaction.Amount, "Transaction amount should be 200")

	t.Log("TestTransferCoinsScenario complete")
}
