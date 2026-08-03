package app

import (
	"github.com/prachii06/LedgerX/internal/config"
	"github.com/prachii06/LedgerX/internal/database"
	"github.com/prachii06/LedgerX/internal/generator"
	"github.com/prachii06/LedgerX/internal/handlers"
	"github.com/prachii06/LedgerX/internal/logger"
	"github.com/prachii06/LedgerX/internal/repository"
	"github.com/prachii06/LedgerX/internal/server"
	"github.com/prachii06/LedgerX/internal/services"
)

func Run() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	// Initialize logger
	log := logger.New()

	log.Info(
		"Starting LedgerX",
		"environment", cfg.AppEnv,
	)

	// Connect to PostgreSQL
	db, err := database.Connect(cfg)
	if err != nil {
		log.Error("Failed to connect to database", "error", err)
		panic(err)
	}
	defer db.Close()

	log.Info("Database connected successfully")

	
	// Dependency Injection
	
	// Repository
	transactionRepo := repository.NewTransactionRepository(db)

	// Services
	transactionService := services.NewTransactionService(transactionRepo)

	// Generators
	transactionGenerator := generator.NewTransactionGenerator(transactionService)

	// Handlers
	transactionHandler := handlers.NewTransactionHandler(transactionService)
	simulationHandler := handlers.NewSimulationHandler(transactionGenerator)

	// Server
	router := server.New(
		cfg.ServerPort,
		db,
		transactionHandler,
		simulationHandler,
	)

	// Start server
	log.Info(
		"HTTP server started",
		"port", cfg.ServerPort,
	)

	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Error("Failed to start HTTP server", "error", err)
		panic(err)
	}
}