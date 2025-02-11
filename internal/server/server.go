package server

import (
	"context"
	"errors"
	"github.com/Egorpalan/avito-shop/internal/handlers"
	"github.com/Egorpalan/avito-shop/internal/middleware"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	router *gin.Engine
}

func NewServer(authHandler *handlers.AuthHandler, infoHandler *handlers.InfoHandler, transactionHandler *handlers.TransactionHandler, merchHandler *handlers.MerchHandler) *Server {
	r := gin.Default()

	// Публичные эндпоинты
	r.POST("/api/auth/register", authHandler.Register)
	r.POST("/api/auth/login", authHandler.Login)

	// Защищенные эндпоинты
	auth := r.Group("/api")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.GET("/info", infoHandler.GetUserInfo)
		auth.POST("/sendCoin", transactionHandler.SendCoins)
		auth.GET("/buy/:item", merchHandler.BuyMerch)
	}

	return &Server{router: r}
}

func (s *Server) Start() {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: s.router,
	}

	// Запуск сервера в отдельной goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Ожидание сигналов для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Создаем контекст с таймаутом для graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Останавливаем сервер
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}

	log.Println("Server stopped")
}
