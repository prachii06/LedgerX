package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prachii06/LedgerX/internal/handlers"
)

func New(port string, db *pgxpool.Pool, transactionHandler *handlers.TransactionHandler, simulationHandler *handlers.SimulationHandler) *gin.Engine {
	router := gin.Default()

	healthHandler := handlers.NewHealthHandler(db)

	router.GET("/health", healthHandler.Health)
	router.GET("/ready", healthHandler.Ready)
	router.GET("/live", healthHandler.Live)
	router.POST("/transactions", transactionHandler.Create)
	router.POST("/simulate", simulationHandler.Generate)
	fmt.Printf("server running on port %s\n", port)

	return router
}
