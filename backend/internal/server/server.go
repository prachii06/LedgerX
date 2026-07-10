package server

import(
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prachii06/LedgerX/internal/handlers"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(port string, db *pgxpool.Pool) *gin.Engine{
	router := gin.Default()

	router.GET("/health",handlers.Health)
	router.GET("/ready", handlers.Ready(db))
	router.GET("/live",handlers.Live)

	fmt.Printf("server running on port %s\n",port)

	return router
}  