package repository

import (
	"errors"

	"github.com/Egorpalan/avito-shop/internal/models"
	"gorm.io/gorm"
)

type TransactionRepositoryInterface interface {
	GetUserBalance(userID uint) (int, error)
	UpdateUserBalance(userID uint, newBalance int) error
	CreateTransaction(transaction *models.Transaction) error
	TransferCoins(fromUserID, toUserID uint, amount int) error
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

func (r *TransactionRepository) TransferCoins(fromUserID, toUserID uint, amount int) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	fromUserBalance, err := r.GetUserBalance(fromUserID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if fromUserBalance < amount {
		tx.Rollback()
		return errors.New("insufficient funds")
	}

	toUserBalance, err := r.GetUserBalance(toUserID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&models.User{}).Where("id = ?", fromUserID).Update("coins", fromUserBalance-amount).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&models.User{}).Where("id = ?", toUserID).Update("coins", toUserBalance+amount).Error; err != nil {
		tx.Rollback()
		return err
	}

	transaction := &models.Transaction{
		FromUserID: fromUserID,
		ToUserID:   toUserID,
		Amount:     amount,
	}
	if err := tx.Create(transaction).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
