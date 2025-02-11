package service

import (
	"github.com/Egorpalan/avito-shop/internal/models"
	"github.com/Egorpalan/avito-shop/internal/repository"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserService_Register(t *testing.T) {
	mockRepo := &repository.MockUserRepository{}
	userService := NewUserService(mockRepo)

	user := &models.User{Username: "testuser", Password: "testpass"}
	mockRepo.On("CreateUser", user).Return(nil)

	err := userService.Register(user)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
