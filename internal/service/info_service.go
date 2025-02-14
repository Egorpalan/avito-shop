package service

import (
	"github.com/Egorpalan/avito-shop/internal/models"
	"github.com/Egorpalan/avito-shop/internal/repository"
	"log"
	"sync"
)

type InfoService struct {
	userRepo *repository.UserRepository
}

func NewInfoService(userRepo *repository.UserRepository) *InfoService {
	return &InfoService{userRepo: userRepo}
}

func (s *InfoService) GetUserInfo(username string) (*models.InfoResponse, error) {
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	var inventory []models.Inventory
	var transactions []models.Transaction
	var wg sync.WaitGroup
	var invErr, transErr error

	wg.Add(2)
	go func() {
		defer wg.Done()
		inventory, invErr = s.userRepo.GetUserInventory(user.ID)
	}()
	go func() {
		defer wg.Done()
		transactions, transErr = s.userRepo.GetUserTransactions(user.ID)
	}()
	wg.Wait()

	if invErr != nil {
		log.Printf("Failed to get user inventory: %v", invErr)
		return nil, invErr
	}
	if transErr != nil {
		log.Printf("Failed to get user transactions: %v", transErr)
		return nil, transErr
	}

	var received []models.TransactionDetail
	var sent []models.TransactionDetail
	for _, t := range transactions {
		if t.ToUserID == user.ID {
			fromUser, err := s.userRepo.GetUserByID(t.FromUserID)
			if err != nil {
				return nil, err
			}
			received = append(received, models.TransactionDetail{
				FromUser: fromUser.Username,
				Amount:   t.Amount,
			})
		} else {
			toUser, err := s.userRepo.GetUserByID(t.ToUserID)
			if err != nil {
				return nil, err
			}
			sent = append(sent, models.TransactionDetail{
				ToUser: toUser.Username,
				Amount: t.Amount,
			})
		}
	}

	return &models.InfoResponse{
		Coins:     user.Coins,
		Inventory: inventory,
		CoinHistory: models.CoinHistory{
			Received: received,
			Sent:     sent,
		},
	}, nil
}
