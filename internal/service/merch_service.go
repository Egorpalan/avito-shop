package service

import (
	"errors"
	"github.com/Egorpalan/avito-shop/internal/models"
	"github.com/Egorpalan/avito-shop/internal/repository"
)

//type MerchService struct {
//	merchRepo       *repository.MerchRepository
//	transactionRepo *repository.TransactionRepository
//	userRepo        *repository.UserRepository
//}
//
//func NewMerchService(merchRepo *repository.MerchRepository, transactionRepo *repository.TransactionRepository, userRepo *repository.UserRepository) *MerchService {
//	return &MerchService{merchRepo: merchRepo, transactionRepo: transactionRepo, userRepo: userRepo}
//}

type MerchService struct {
	merchRepo       repository.MerchRepositoryInterface
	transactionRepo repository.TransactionRepositoryInterface
	userRepo        repository.UserRepositoryInterface
}

func NewMerchService(
	merchRepo repository.MerchRepositoryInterface,
	transactionRepo repository.TransactionRepositoryInterface,
	userRepo repository.UserRepositoryInterface,
) *MerchService {
	return &MerchService{
		merchRepo:       merchRepo,
		transactionRepo: transactionRepo,
		userRepo:        userRepo,
	}
}

func (s *MerchService) BuyMerch(userID, merchID uint) error {
	price, err := s.merchRepo.GetMerchPrice(merchID)
	if err != nil {
		return err
	}

	userBalance, err := s.transactionRepo.GetUserBalance(userID)
	if err != nil {
		return err
	}

	if userBalance < price {
		return errors.New("insufficient funds")
	}

	if err := s.transactionRepo.UpdateUserBalance(userID, userBalance-price); err != nil {
		return err
	}

	return s.merchRepo.AddToInventory(userID, merchID, 1)
}

func (s *MerchService) GetMerchByName(name string) (*models.Merch, error) {
	return s.merchRepo.GetMerchByName(name)
}

func (s *MerchService) BuyMerchByUsername(username string, merchID uint) error {
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return errors.New("user not found")
	}

	return s.BuyMerch(user.ID, merchID)
}
