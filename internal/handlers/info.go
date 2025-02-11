package handlers

import (
	"github.com/Egorpalan/avito-shop/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type InfoHandler struct {
	infoService *service.InfoService
}

func NewInfoHandler(infoService *service.InfoService) *InfoHandler {
	return &InfoHandler{infoService: infoService}
}

func (h *InfoHandler) GetUserInfo(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"errors": "User not authenticated"})
		return
	}

	info, err := h.infoService.GetUserInfo(username.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": err.Error()})
		return
	}

	c.JSON(http.StatusOK, info)
}
