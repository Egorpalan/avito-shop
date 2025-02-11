package service_test

import (
	"errors"
	"github.com/Egorpalan/avito-shop/internal/models"
	"gorm.io/gorm"
	"testing"

	"github.com/Egorpalan/avito-shop/internal/service"
	mockRepo "github.com/Egorpalan/avito-shop/tests/mocks/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTransactionService_SendCoins_Success(t *testing.T) {
	mockTransactionRepo := new(mockRepo.MockTransactionRepository)
	mockUserRepo := new(mockRepo.MockUserRepository)

	transactionService := service.NewTransactionService(mockTransactionRepo, mockUserRepo)

	mockTransactionRepo.On("GetUserBalance", uint(1)).Return(500, nil)
	mockTransactionRepo.On("GetUserBalance", uint(2)).Return(200, nil)
	mockTransactionRepo.On("UpdateUserBalance", uint(1), 300).Return(nil)
	mockTransactionRepo.On("UpdateUserBalance", uint(2), 400).Return(nil)
	mockTransactionRepo.On("CreateTransaction", mock.Anything).Return(nil)

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

func TestTransactionService_SendCoins_InsufficientFunds(t *testing.T) {
	mockTransactionRepo := new(mockRepo.MockTransactionRepository)
	mockUserRepo := new(mockRepo.MockUserRepository)

	transactionService := service.NewTransactionService(mockTransactionRepo, mockUserRepo)

	mockTransactionRepo.On("GetUserBalance", uint(1)).Return(50, nil) // Баланс меньше, чем отправляем

	err := transactionService.SendCoins(1, 2, 100)

	assert.Error(t, err)
	assert.Equal(t, "insufficient funds", err.Error())

	mockTransactionRepo.AssertExpectations(t)
}

func TestTransactionService_SendCoins_UserNotFound(t *testing.T) {
	mockTransactionRepo := new(mockRepo.MockTransactionRepository)
	mockUserRepo := new(mockRepo.MockUserRepository)

	transactionService := service.NewTransactionService(mockTransactionRepo, mockUserRepo)

	mockTransactionRepo.On("GetUserBalance", uint(1)).Return(500, nil)
	mockTransactionRepo.On("GetUserBalance", uint(2)).Return(0, errors.New("user not found")) // Ошибка на втором юзере

	err := transactionService.SendCoins(1, 2, 100)

	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())

	mockTransactionRepo.AssertExpectations(t)
}

func TestTransactionService_SendCoinsByUsername_Success(t *testing.T) {
	mockTransactionRepo := new(mockRepo.MockTransactionRepository)
	mockUserRepo := new(mockRepo.MockUserRepository)

	transactionService := service.NewTransactionService(mockTransactionRepo, mockUserRepo)

	mockUserRepo.On("GetUserByUsername", "alice").Return(&models.User{Model: gorm.Model{ID: 1}, Username: "alice"}, nil)
	mockUserRepo.On("GetUserByUsername", "bob").Return(&models.User{Model: gorm.Model{ID: 2}, Username: "bob"}, nil)

	mockTransactionRepo.On("GetUserBalance", uint(1)).Return(500, nil)
	mockTransactionRepo.On("GetUserBalance", uint(2)).Return(100, nil)
	mockTransactionRepo.On("UpdateUserBalance", uint(1), 300).Return(nil)
	mockTransactionRepo.On("UpdateUserBalance", uint(2), 300).Return(nil)
	mockTransactionRepo.On("CreateTransaction", mock.Anything).Return(nil)

	err := transactionService.SendCoinsByUsername("alice", "bob", 200)

	assert.NoError(t, err)

	mockUserRepo.AssertExpectations(t)
	mockTransactionRepo.AssertExpectations(t)
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
