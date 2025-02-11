package repository

import (
	"github.com/Egorpalan/avito-shop/internal/models"
	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	CreateUser(user *models.User) error
	GetUserByUsername(username string) (*models.User, error)
	GetUserByID(id uint) (*models.User, error)
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserInventory(userID uint) ([]models.Inventory, error) {
	var inventory []models.Inventory
	if err := r.db.Where("user_id = ?", userID).Find(&inventory).Error; err != nil {
		return nil, err
	}
	return inventory, nil
}

func (r *UserRepository) GetUserTransactions(userID uint) ([]models.Transaction, error) {
	var transactions []models.Transaction
	if err := r.db.Where("from_user_id = ? OR to_user_id = ?", userID, userID).Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}
