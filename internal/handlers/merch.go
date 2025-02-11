package handlers

import (
	"github.com/Egorpalan/avito-shop/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type MerchHandler struct {
	merchService *service.MerchService
}

func NewMerchHandler(merchService *service.MerchService) *MerchHandler {
	return &MerchHandler{merchService: merchService}
}

func (h *MerchHandler) BuyMerch(c *gin.Context) {
	// Получаем username из контекста
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"errors": "User not authenticated"})
		return
	}

	// Преобразуем username в строку
	fromUsername, ok := username.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": "Invalid username type"})
		return
	}

	item := c.Param("item")
	if item == "" {
		c.JSON(http.StatusBadRequest, gin.H{"errors": "Invalid item"})
		return
	}

	// Найдем товар по имени через сервис
	merch, err := h.merchService.GetMerchByName(item)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": "Item not found"})
		return
	}

	// Передаем username в сервис
	if err := h.merchService.BuyMerchByUsername(fromUsername, merch.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
