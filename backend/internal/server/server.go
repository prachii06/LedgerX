package server

import(
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prachii06/LedgerX/internal/handlers"
)

func New(port string) *gin.Engine{
	router := gin.Default()

	router.GET("/health",handlers.Health)
	router.GET("/ready",handlers.Ready)
	router.GET("/live",handlers.Live)

	fmt.Printf("server running on port %s\n",port)

	return router
}  