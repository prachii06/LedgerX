package app

import (
	"github.com/prachii06/LedgerX/internal/config"
	"github.com/prachii06/LedgerX/internal/logger"
	"github.com/prachii06/LedgerX/internal/server"
	"github.com/prachii06/LedgerX/internal/database"
	
)

func Run() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	log := logger.New()

	log.Info(
		"Starting LedgerX",
		"environment", cfg.AppEnv,
	)

	db, err := database.Connect(cfg)
	if err != nil {
	log.Error("Failed to connect to database", "error", err)
	panic(err)
}
	defer db.Close()

log.Info("Database connected successfully")

	log.Info("db connected successfully")

	router := server.New(cfg.ServerPort, db)
	
	err = router.Run(":" + cfg.ServerPort)
	if err != nil {
		log.Error("Failed to start server", "error", err)
		panic(err)
	}
}