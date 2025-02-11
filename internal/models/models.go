package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Coins    int    `gorm:"default:1000"`
}

type Merch struct {
	gorm.Model
	Name  string `gorm:"unique;not null"`
	Price int    `gorm:"not null"`
}

type Transaction struct {
	gorm.Model
	FromUserID uint
	ToUserID   uint
	Amount     int
}

type Inventory struct {
	gorm.Model
	UserID   uint
	MerchID  uint
	Quantity int
}

type InfoResponse struct {
	Coins       int         `json:"coins"`
	Inventory   []Inventory `json:"inventory"`
	CoinHistory CoinHistory `json:"coinHistory"`
}

type CoinHistory struct {
	Received []TransactionDetail `json:"received"`
	Sent     []TransactionDetail `json:"sent"`
}

type TransactionDetail struct {
	FromUser string `json:"fromUser,omitempty"`
	ToUser   string `json:"toUser,omitempty"`
	Amount   int    `json:"amount"`
}
