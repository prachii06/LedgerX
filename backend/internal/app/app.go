package app

import (
	"github.com/prachii06/LedgerX/internal/config"
	"github.com/prachii06/LedgerX/internal/logger"
	"github.com/prachii06/LedgerX/internal/server"
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
     
	router := server.New(cfg.ServerPort)

	err = router.Run(":" + cfg.ServerPort)
    if err != nil {
		log.Error("Failed to start server", "error", err)
		panic(err)
}
	
}