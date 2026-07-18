package server

import(
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prachii06/LedgerX/internal/handlers"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prachii06/LedgerX/internal/services"
)

func New(port string, db *pgxpool.Pool) *gin.Engine{
	router := gin.Default()

	healthHandler := handlers.NewHealthHandler(db)
	
	router.GET("/health", healthHandler.Health)
	router.GET("/ready", healthHandler.Ready)
	router.GET("/live", healthHandler.Live)
	router.POST("/transactions", handlers.NewTransactionHandler(db).CreateTransaction)
	
	fmt.Printf("server running on port %s\n",port)

	return router
}  