package main

import (
	"github.com/Egorpalan/avito-shop/config"
	"github.com/Egorpalan/avito-shop/internal/handlers"
	"github.com/Egorpalan/avito-shop/internal/repository"
	"github.com/Egorpalan/avito-shop/internal/server"
	"github.com/Egorpalan/avito-shop/internal/service"
	"github.com/Egorpalan/avito-shop/migrations"
	"github.com/Egorpalan/avito-shop/pkg/db"
	"log"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.LoadConfig(".env.example")

	// Инициализация базы данных
	dbConn, err := db.InitDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Выполняем миграции
	err = migrations.AutoMigrate(dbConn)
	if err != nil {
		log.Fatal("Failed to auto migrate:", err)
	}

	// Инициализация репозиториев
	userRepo := repository.NewUserRepository(dbConn)
	transactionRepo := repository.NewTransactionRepository(dbConn)
	merchRepo := repository.NewMerchRepository(dbConn)

	// Инициализация сервисов
	userService := service.NewUserService(userRepo)
	transactionService := service.NewTransactionService(transactionRepo, userRepo)
	merchService := service.NewMerchService(merchRepo, transactionRepo, userRepo)
	infoService := service.NewInfoService(userRepo)

	// Инициализация хендлеров
	authHandler := handlers.NewAuthHandler(userService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)
	merchHandler := handlers.NewMerchHandler(merchService)
	infoHandler := handlers.NewInfoHandler(infoService)

	// Создаем и запускаем сервер
	server := server.NewServer(authHandler, infoHandler, transactionHandler, merchHandler)
	server.Start()
}
