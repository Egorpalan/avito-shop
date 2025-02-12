package main

import (
	"github.com/Egorpalan/avito-shop/config"
	"github.com/Egorpalan/avito-shop/internal/handlers"
	"github.com/Egorpalan/avito-shop/internal/repository"
	"github.com/Egorpalan/avito-shop/internal/server"
	"github.com/Egorpalan/avito-shop/internal/service"
	"github.com/Egorpalan/avito-shop/migrations"
	"github.com/Egorpalan/avito-shop/pkg/db"
	"github.com/Egorpalan/avito-shop/pkg/logger"
)

func main() {
	log := logger.InitLogger()
	cfg := config.LoadConfig(".env.example")
	log.Info("Config loaded", "env", "env.example")

	dbConn, err := db.InitDB(cfg)
	if err != nil {
		log.Error("Failed to connect to database", "error", err)
	}
	log.Info("Connected to database")

	err = migrations.AutoMigrate(dbConn)
	if err := migrations.AutoMigrate(dbConn); err != nil {
		log.Error("Failed to auto migrate", "error", err)
	}
	log.Info("Database migrated successfully")

	userRepo := repository.NewUserRepository(dbConn)
	transactionRepo := repository.NewTransactionRepository(dbConn)
	merchRepo := repository.NewMerchRepository(dbConn)

	userService := service.NewUserService(userRepo)
	transactionService := service.NewTransactionService(transactionRepo, userRepo)
	merchService := service.NewMerchService(merchRepo, transactionRepo, userRepo)
	infoService := service.NewInfoService(userRepo)

	authHandler := handlers.NewAuthHandler(userService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)
	merchHandler := handlers.NewMerchHandler(merchService)
	infoHandler := handlers.NewInfoHandler(infoService)

	server := server.NewServer(authHandler, infoHandler, transactionHandler, merchHandler)
	server.Start()
}
