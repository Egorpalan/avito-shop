package handlers

import (
	"github.com/Egorpalan/avito-shop/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TransactionHandler struct {
	transactionService *service.TransactionService
}

func NewTransactionHandler(transactionService *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{transactionService: transactionService}
}

func (h *TransactionHandler) SendCoins(c *gin.Context) {
	var request struct {
		ToUser string `json:"toUser"`
		Amount int    `json:"amount"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": "Invalid request"})
		return
	}

	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"errors": "User not authenticated"})
		return
	}

	fromUsername, ok := username.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": "Invalid username type"})
		return
	}

	if err := h.transactionService.SendCoinsByUsername(fromUsername, request.ToUser, request.Amount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
