package service

import (
	"errors"
	"github.com/Egorpalan/avito-shop/internal/models"
	"github.com/Egorpalan/avito-shop/internal/repository"
)

type TransactionService struct {
	transactionRepo *repository.TransactionRepository
	userRepo        *repository.UserRepository
}

func NewTransactionService(transactionRepo *repository.TransactionRepository, userRepo *repository.UserRepository) *TransactionService {
	return &TransactionService{transactionRepo: transactionRepo, userRepo: userRepo}
}

func (s *TransactionService) SendCoins(fromUserID, toUserID uint, amount int) error {
	if fromUserID == toUserID {
		return errors.New("cannot send coins to yourself")
	}

	fromUserBalance, err := s.transactionRepo.GetUserBalance(fromUserID)
	if err != nil {
		return err
	}

	if fromUserBalance < amount {
		return errors.New("insufficient funds")
	}

	toUserBalance, err := s.transactionRepo.GetUserBalance(toUserID)
	if err != nil {
		return err
	}

	// Обновляем балансы
	if err := s.transactionRepo.UpdateUserBalance(fromUserID, fromUserBalance-amount); err != nil {
		return err
	}
	if err := s.transactionRepo.UpdateUserBalance(toUserID, toUserBalance+amount); err != nil {
		return err
	}

	// Создаем запись о транзакции
	transaction := &models.Transaction{
		FromUserID: fromUserID,
		ToUserID:   toUserID,
		Amount:     amount,
	}
	return s.transactionRepo.CreateTransaction(transaction)
}

func (s *TransactionService) SendCoinsByUsername(fromUsername, toUsername string, amount int) error {
	// Находим fromUser по username
	fromUser, err := s.userRepo.GetUserByUsername(fromUsername)
	if err != nil {
		return errors.New("sender user not found")
	}

	// Находим toUser по username
	toUser, err := s.userRepo.GetUserByUsername(toUsername)
	if err != nil {
		return errors.New("recipient user not found")
	}

	// Выполняем передачу монет
	return s.SendCoins(fromUser.ID, toUser.ID, amount)
}

func (s *TransactionService) GetUserByUsername(username string) (*models.User, error) {
	return s.userRepo.GetUserByUsername(username)
}
