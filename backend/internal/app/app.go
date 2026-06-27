package app

import (
	"github.com/prachii06/LedgerX/internal/config"
	"github.com/prachii06/LedgerX/internal/logger"
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
}