package repository

import (
	"github.com/Egorpalan/avito-shop/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) GetUserBalance(userID uint) (int, error) {
	args := m.Called(userID)
	return args.Int(0), args.Error(1)
}

func (m *MockTransactionRepository) UpdateUserBalance(userID uint, newBalance int) error {
	args := m.Called(userID, newBalance)
	return args.Error(0)
}

func (m *MockTransactionRepository) CreateTransaction(transaction *models.Transaction) error {
	args := m.Called(transaction)
	return args.Error(0)
}

func (m *MockTransactionRepository) TransferCoins(fromUserID, toUserID uint, amount int) error {
	args := m.Called(fromUserID, toUserID, amount)
	return args.Error(0)
}
