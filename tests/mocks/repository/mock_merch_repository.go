package repository

import (
	"github.com/Egorpalan/avito-shop/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockMerchRepository struct {
	mock.Mock
}

func (m *MockMerchRepository) GetMerchPrice(merchID uint) (int, error) {
	args := m.Called(merchID)
	return args.Int(0), args.Error(1)
}

func (m *MockMerchRepository) AddToInventory(userID, merchID uint, quantity int) error {
	args := m.Called(userID, merchID, quantity)
	return args.Error(0)
}

func (m *MockMerchRepository) GetMerchByName(name string) (*models.Merch, error) {
	args := m.Called(name)
	return args.Get(0).(*models.Merch), args.Error(1)
}
