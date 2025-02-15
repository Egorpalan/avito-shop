package server

import (
	"context"
	"errors"
	"github.com/Egorpalan/avito-shop/config"
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

func NewServer(authHandler *handlers.AuthHandler, infoHandler *handlers.InfoHandler, transactionHandler *handlers.TransactionHandler, merchHandler *handlers.MerchHandler, cfg *config.Config) *Server {
	r := gin.Default()

	r.POST("/api/auth/register", authHandler.Register)
	r.POST("/api/auth/login", authHandler.Login)

	auth := r.Group("/api")
	auth.Use(middleware.AuthMiddleware(cfg))
	{
		auth.GET("/info", infoHandler.GetUserInfo)
		auth.POST("/sendCoin", transactionHandler.SendCoins)
		auth.GET("/buy/:item", merchHandler.BuyMerch)
	}

	return &Server{router: r}
}

func (s *Server) Start() {
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      s.router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}

	log.Println("Server stopped")
}
