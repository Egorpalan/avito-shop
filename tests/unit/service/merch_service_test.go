package service_test

import (
	"errors"
	"gorm.io/gorm"
	"testing"

	"github.com/Egorpalan/avito-shop/internal/models"
	"github.com/Egorpalan/avito-shop/internal/service"
	mockRepo "github.com/Egorpalan/avito-shop/tests/mocks/repository"
	"github.com/stretchr/testify/assert"
)

func TestMerchService_BuyMerch_Success(t *testing.T) {
	mockMerchRepo := new(mockRepo.MockMerchRepository)
	mockTransactionRepo := new(mockRepo.MockTransactionRepository)
	mockUserRepo := new(mockRepo.MockUserRepository)

	merchService := service.NewMerchService(mockMerchRepo, mockTransactionRepo, mockUserRepo)

	userID := uint(1)
	merchID := uint(2)
	merchPrice := 100
	userBalance := 200

	mockMerchRepo.On("GetMerchPrice", merchID).Return(merchPrice, nil)
	mockTransactionRepo.On("GetUserBalance", userID).Return(userBalance, nil)
	mockTransactionRepo.On("UpdateUserBalance", userID, userBalance-merchPrice).Return(nil)
	mockMerchRepo.On("AddToInventory", userID, merchID, 1).Return(nil)

	err := merchService.BuyMerch(userID, merchID)

	assert.NoError(t, err)
	mockMerchRepo.AssertExpectations(t)
	mockTransactionRepo.AssertExpectations(t)
}

func TestMerchService_BuyMerch_InsufficientFunds(t *testing.T) {
	mockMerchRepo := new(mockRepo.MockMerchRepository)
	mockTransactionRepo := new(mockRepo.MockTransactionRepository)
	mockUserRepo := new(mockRepo.MockUserRepository)

	merchService := service.NewMerchService(mockMerchRepo, mockTransactionRepo, mockUserRepo)

	userID := uint(1)
	merchID := uint(2)
	merchPrice := 100
	userBalance := 50

	mockMerchRepo.On("GetMerchPrice", merchID).Return(merchPrice, nil)
	mockTransactionRepo.On("GetUserBalance", userID).Return(userBalance, nil)

	err := merchService.BuyMerch(userID, merchID)

	assert.Error(t, err)
	assert.Equal(t, "insufficient funds", err.Error())
	mockMerchRepo.AssertExpectations(t)
	mockTransactionRepo.AssertExpectations(t)
}

func TestMerchService_GetMerchByName_Success(t *testing.T) {
	mockMerchRepo := new(mockRepo.MockMerchRepository)
	mockTransactionRepo := new(mockRepo.MockTransactionRepository)
	mockUserRepo := new(mockRepo.MockUserRepository)

	merchService := service.NewMerchService(mockMerchRepo, mockTransactionRepo, mockUserRepo)

	expectedMerch := &models.Merch{Name: "T-shirt", Price: 500}
	mockMerchRepo.On("GetMerchByName", "T-shirt").Return(expectedMerch, nil)

	merch, err := merchService.GetMerchByName("T-shirt")

	assert.NoError(t, err)
	assert.Equal(t, expectedMerch, merch)
	mockMerchRepo.AssertExpectations(t)
}

func TestMerchService_BuyMerchByUsername_Success(t *testing.T) {
	mockMerchRepo := new(mockRepo.MockMerchRepository)
	mockTransactionRepo := new(mockRepo.MockTransactionRepository)
	mockUserRepo := new(mockRepo.MockUserRepository)

	merchService := service.NewMerchService(mockMerchRepo, mockTransactionRepo, mockUserRepo)

	username := "testuser"
	userID := uint(1)
	merchID := uint(2)
	merchPrice := 100
	userBalance := 200

	mockUserRepo.On("GetUserByUsername", "testuser").Return(&models.User{
		Model:    gorm.Model{ID: 1},
		Username: "testuser",
		Password: "hashedpassword",
		Coins:    1000,
	}, nil)
	mockMerchRepo.On("GetMerchPrice", merchID).Return(merchPrice, nil)
	mockTransactionRepo.On("GetUserBalance", userID).Return(userBalance, nil)
	mockTransactionRepo.On("UpdateUserBalance", userID, userBalance-merchPrice).Return(nil)
	mockMerchRepo.On("AddToInventory", userID, merchID, 1).Return(nil)

	err := merchService.BuyMerchByUsername(username, merchID)

	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
	mockMerchRepo.AssertExpectations(t)
	mockTransactionRepo.AssertExpectations(t)
}

func TestMerchService_BuyMerchByUsername_UserNotFound(t *testing.T) {
	mockMerchRepo := new(mockRepo.MockMerchRepository)
	mockTransactionRepo := new(mockRepo.MockTransactionRepository)
	mockUserRepo := new(mockRepo.MockUserRepository)

	merchService := service.NewMerchService(mockMerchRepo, mockTransactionRepo, mockUserRepo)

	username := "unknownuser"
	merchID := uint(2)

	mockUserRepo.On("GetUserByUsername", username).Return(nil, errors.New("user not found"))

	err := merchService.BuyMerchByUsername(username, merchID)

	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
	mockUserRepo.AssertExpectations(t)
}
