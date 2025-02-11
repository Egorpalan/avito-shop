package repository

import (
	"github.com/Egorpalan/avito-shop/internal/models"
	"gorm.io/gorm"
)

type TransactionRepositoryInterface interface {
	GetUserBalance(userID uint) (int, error)
	UpdateUserBalance(userID uint, newBalance int) error
	CreateTransaction(transaction *models.Transaction) error
}

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) CreateTransaction(transaction *models.Transaction) error {
	return r.db.Create(transaction).Error
}

func (r *TransactionRepository) GetUserBalance(userID uint) (int, error) {
	var user models.User
	if err := r.db.First(&user, userID).Error; err != nil {
		return 0, err
	}
	return user.Coins, nil
}

func (r *TransactionRepository) UpdateUserBalance(userID uint, newBalance int) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).Update("coins", newBalance).Error
}
