package service

import (
	"github.com/Egorpalan/avito-shop/internal/models"
	"github.com/Egorpalan/avito-shop/internal/repository"
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

	inventory, err := s.userRepo.GetUserInventory(user.ID)
	if err != nil {
		return nil, err
	}

	transactions, err := s.userRepo.GetUserTransactions(user.ID)
	if err != nil {
		return nil, err
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
