package service_test

import (
	"golang.org/x/crypto/bcrypt"
	"testing"

	"github.com/Egorpalan/avito-shop/internal/models"
	"github.com/Egorpalan/avito-shop/internal/service"
	mockRepo "github.com/Egorpalan/avito-shop/tests/mocks/repository"
	"github.com/stretchr/testify/assert"
)

func TestUserService_Register(t *testing.T) {
	mockUserRepo := new(mockRepo.MockUserRepository)

	userService := service.NewUserService(mockUserRepo)

	user := &models.User{Username: "testuser", Password: "testpass"}
	mockUserRepo.On("CreateUser", user).Return(nil)

	err := userService.Register(user)

	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
}

func TestUserService_Login(t *testing.T) {
	mockUserRepo := new(mockRepo.MockUserRepository)
	userService := service.NewUserService(mockUserRepo)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testpass"), bcrypt.DefaultCost)

	user := &models.User{Username: "testuser", Password: string(hashedPassword)}
	mockUserRepo.On("GetUserByUsername", "testuser").Return(user, nil)

	_, err := userService.Login("testuser", "testpass")

	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
}
