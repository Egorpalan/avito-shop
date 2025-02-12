package service_test

import (
	"errors"
	"github.com/Egorpalan/avito-shop/internal/models"
	"gorm.io/gorm"
	"testing"

	"github.com/Egorpalan/avito-shop/internal/service"
	mockRepo "github.com/Egorpalan/avito-shop/tests/mocks/repository"
	"github.com/stretchr/testify/assert"
)

func TestTransactionService_SendCoins_Success(t *testing.T) {
	mockTransactionRepo := new(mockRepo.MockTransactionRepository)
	mockUserRepo := new(mockRepo.MockUserRepository)

	transactionService := service.NewTransactionService(mockTransactionRepo, mockUserRepo)

	mockTransactionRepo.On("TransferCoins", uint(1), uint(2), 200).Return(nil)

	err := transactionService.SendCoins(1, 2, 200)

	assert.NoError(t, err)
	mockTransactionRepo.AssertExpectations(t)
}

func TestTransactionService_SendCoins_ToSelfError(t *testing.T) {
	mockTransactionRepo := new(mockRepo.MockTransactionRepository)
	mockUserRepo := new(mockRepo.MockUserRepository)

	transactionService := service.NewTransactionService(mockTransactionRepo, mockUserRepo)

	err := transactionService.SendCoins(1, 1, 100)

	assert.Error(t, err)
	assert.Equal(t, "cannot send coins to yourself", err.Error())

	mockTransactionRepo.AssertNotCalled(t, "GetUserBalance")
}

func TestTransactionService_SendCoinsByUsername_SenderNotFound(t *testing.T) {
	mockTransactionRepo := new(mockRepo.MockTransactionRepository)
	mockUserRepo := new(mockRepo.MockUserRepository)

	transactionService := service.NewTransactionService(mockTransactionRepo, mockUserRepo)

	mockUserRepo.On("GetUserByUsername", "alice").Return(nil, errors.New("sender user not found"))

	err := transactionService.SendCoinsByUsername("alice", "bob", 100)

	assert.Error(t, err)
	assert.Equal(t, "sender user not found", err.Error())

	mockUserRepo.AssertExpectations(t)
}

func TestTransactionService_SendCoinsByUsername_RecipientNotFound(t *testing.T) {
	mockTransactionRepo := new(mockRepo.MockTransactionRepository)
	mockUserRepo := new(mockRepo.MockUserRepository)

	transactionService := service.NewTransactionService(mockTransactionRepo, mockUserRepo)

	mockUserRepo.On("GetUserByUsername", "alice").Return(&models.User{Model: gorm.Model{ID: 1}, Username: "alice"}, nil)
	mockUserRepo.On("GetUserByUsername", "bob").Return(nil, errors.New("recipient user not found"))

	err := transactionService.SendCoinsByUsername("alice", "bob", 100)

	assert.Error(t, err)
	assert.Equal(t, "recipient user not found", err.Error())

	mockUserRepo.AssertExpectations(t)
}
