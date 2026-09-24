package app

import (
	"context"
	"time"

	"github.com/prachii06/LedgerX/internal/cache"
	"github.com/prachii06/LedgerX/internal/config"
	"github.com/prachii06/LedgerX/internal/database"
	"github.com/prachii06/LedgerX/internal/generator"
	"github.com/prachii06/LedgerX/internal/handlers"
	"github.com/prachii06/LedgerX/internal/kafka"
	"github.com/prachii06/LedgerX/internal/logger"
	"github.com/prachii06/LedgerX/internal/metrics"
	"github.com/prachii06/LedgerX/internal/redis"
	"github.com/prachii06/LedgerX/internal/repository"
	"github.com/prachii06/LedgerX/internal/server"
	"github.com/prachii06/LedgerX/internal/services"
	"github.com/prachii06/LedgerX/internal/worker"
)

func Run() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	// Initialize metrics
	metrics.Init()

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

	// Connect to Redis
	redisClient := redis.NewClient(cfg)
	defer redisClient.Close()

	if err := redisClient.Ping(context.Background()); err != nil {
		log.Error("Failed to connect to Redis", "error", err)
		panic(err)
	}

	log.Info("Redis connected successfully")

	// ----------------------------
	// Dependency Injection
	// ----------------------------

	// Repositories
	transactionRepo := repository.NewTransactionRepository(db)
	reconciliationRepo := repository.NewReconciliationRepository(db)
	reconciliationResultRepo := repository.NewReconciliationResultRepository(db)
	eventRepo := repository.NewEventRepository(db)

	// Cache
	reconciliationCache := cache.NewReconciliationCache(redisClient.Client, 10*time.Minute)

	// Services
	transactionService := services.NewTransactionService(transactionRepo)
	reconciliationService := services.NewReconciliationService(
		reconciliationRepo,
		reconciliationResultRepo,
		reconciliationCache,
	)
	eventService := services.NewEventService(eventRepo)

	// Kafka
	kafkaProducer := kafka.NewProducer(cfg.KafkaBrokers)
	defer kafkaProducer.Close()

	kafkaConsumer := kafka.NewConsumer(
		cfg.KafkaBrokers,
		"ledgerx-event-consumer",
		eventService,
		kafkaProducer,
	)
	defer kafkaConsumer.Close()

	go kafkaConsumer.Start(context.Background()) //start kakfa consumer

	// Background Workers
	reconciliationWorker := worker.NewReconciliationWorker(
		reconciliationRepo,
		reconciliationService,
	)

	go reconciliationWorker.Start(
		context.Background(),
		5*time.Second,
	)

	// Generators
	transactionGenerator := generator.NewTransactionGenerator(transactionService)
	eventGenerator := generator.NewEventGenerator(kafkaProducer)

	// Application Services
	simulationService := services.NewSimulationService(
		transactionGenerator,
		eventGenerator,
	)

	// Handlers
	transactionHandler := handlers.NewTransactionHandler(transactionService)
	simulationHandler := handlers.NewSimulationHandler(simulationService)
	reconciliationHandler := handlers.NewReconciliationHandler(reconciliationService)
	eventHandler := handlers.NewEventHandler(eventService)
	dashboardHandler := handlers.NewDashboardHandler(transactionService)

	// Server
	router := server.New(
		cfg.ServerPort,
		cfg.CORSAllowedOrigins,
		db,
		transactionHandler,
		simulationHandler,
		reconciliationHandler,
		eventHandler,
		dashboardHandler,
	)


	// Start HTTP Server
	log.Info(
		"HTTP server started",
		"port", cfg.ServerPort,
	)

	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Error("Failed to start HTTP server", "error", err)
		panic(err)
	}
}
