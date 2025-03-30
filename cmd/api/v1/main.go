package main

import (
	"fmt"
	"log"
	"time"

	"lumon-backend/internal/config"
	"lumon-backend/internal/handler"
	"lumon-backend/internal/migrations"
	"lumon-backend/internal/repository/database"
	"lumon-backend/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	result := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	if result.Error != nil {
		log.Fatal("failed to enable uuid-ossp extension: %w", result.Error)
	}

	err = db.AutoMigrate(migrations.GetMigrationModels()...)
	if err != nil {
		log.Fatal("Failed to perform Database Migrations")
	}

	userRepo := database.NewUserRepository(db)
	documentRepo := database.NewDocumentRepository(db)
	transactionRepo := database.NewTransactionRepository(db)
	loanRequestRepo := database.NewLoanRequestRepository(db)
	accountRepo := database.NewAccountRepository(db)
	walletRepo := database.NewWalletRepository(db)

	documentService := service.NewDocumentService(documentRepo)
	transactionService := service.NewTransactionService(transactionRepo)
	userService := service.NewUserService(userRepo)
	loanRequestService := service.NewLoanRequestService(loanRequestRepo)
	accountService := service.NewAccountService(accountRepo)
	walletService := service.NewWalletService(walletRepo)

	userHandler := handler.NewUserHandler(userService, cfg)
	authHandler := handler.NewAuthHandler(userService, cfg)
	documentHandler := handler.NewDocumentHandler(documentService, cfg)
	transactionHandler := handler.NewTransactionHandler(transactionService, userService, cfg)
	chatHandler := handler.NewChatHandler(transactionService, cfg)
	loanRequestHandler := handler.NewLoanRequestHandler(loanRequestService, userService, cfg)
	accountHandler := handler.NewAccountHandler(accountService, cfg)
	walletHandler := handler.NewWalletsHandler(walletService, cfg)

	r := gin.Default()

  r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", 
		"https://lumon-tech-dashboard.vercel.app", 
		"https://lumon-tech-dashboard-fawn.vercel.app",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge: 12 * time.Hour,
	}))

	api := r.Group("/api")
	{
		userHandler.RegisterRoutes(api)
		authHandler.RegisterRoutes(api)
		chatHandler.RegisterRoutes(api)
		documentHandler.RegisterRoutes(api)
		transactionHandler.RegisterRoutes(api)
		loanRequestHandler.RegisterRoutes(api)
		accountHandler.RegisterRoutes(api)
		walletHandler.RegisterRoutes(api)
	}

	if err := r.Run(fmt.Sprintf(":%s", cfg.Port)); err != nil {
		log.Fatal("Failed to Start Server, Gracefully Stopping")
	}
}
