package repository

import (
	"github.com/Egorpalan/avito-shop/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	user, _ := args.Get(0).(*models.User)
	return user, args.Error(1)
}

func (m *MockUserRepository) GetUserByID(id uint) (*models.User, error) {
	args := m.Called(id)
	user, _ := args.Get(0).(*models.User)
	return user, args.Error(1)
}
