package repository

import (
	"github.com/Egorpalan/avito-shop/internal/models"
	"gorm.io/gorm"
)

type MerchRepository struct {
	db *gorm.DB
}

func NewMerchRepository(db *gorm.DB) *MerchRepository {
	return &MerchRepository{db: db}
}

func (r *MerchRepository) GetMerchPrice(merchID uint) (int, error) {
	var merch models.Merch
	if err := r.db.First(&merch, merchID).Error; err != nil {
		return 0, err
	}
	return merch.Price, nil
}

func (r *MerchRepository) AddToInventory(userID, merchID uint, quantity int) error {
	inventory := &models.Inventory{
		UserID:   userID,
		MerchID:  merchID,
		Quantity: quantity,
	}
	return r.db.Create(inventory).Error
}

func (r *MerchRepository) GetMerchByName(name string) (*models.Merch, error) {
	var merch models.Merch
	if err := r.db.Where("name = ?", name).First(&merch).Error; err != nil {
		return nil, err
	}
	return &merch, nil
}
