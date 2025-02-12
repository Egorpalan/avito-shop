package service

import (
	"errors"

	"github.com/Egorpalan/avito-shop/internal/models"
	"github.com/Egorpalan/avito-shop/internal/repository"
)

type TransactionService struct {
	transactionRepo repository.TransactionRepositoryInterface
	userRepo        repository.UserRepositoryInterface
}

func NewTransactionService(transactionRepo repository.TransactionRepositoryInterface, userRepo repository.UserRepositoryInterface) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		userRepo:        userRepo,
	}
}

// Отправка монет между пользователями
func (s *TransactionService) SendCoins(fromUserID, toUserID uint, amount int) error {
	if fromUserID == toUserID {
		return errors.New("cannot send coins to yourself")
	}

	if amount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	return s.transactionRepo.TransferCoins(fromUserID, toUserID, amount)
}

func (s *TransactionService) SendCoinsByUsername(fromUsername, toUsername string, amount int) error {
	fromUser, err := s.userRepo.GetUserByUsername(fromUsername)
	if err != nil {
		return errors.New("sender user not found")
	}

	toUser, err := s.userRepo.GetUserByUsername(toUsername)
	if err != nil {
		return errors.New("recipient user not found")
	}

	return s.SendCoins(fromUser.ID, toUser.ID, amount)
}

func (s *TransactionService) GetUserByUsername(username string) (*models.User, error) {
	return s.userRepo.GetUserByUsername(username)
}
