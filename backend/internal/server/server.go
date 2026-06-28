package server

import(
	"fmt"
	"github.com/gin-gonic/gin"
)

func New(port string) *gin.Engine{
	router := gin.Default()

	router.GET("/",func(c *gin.Context){
		c.JSON(200, gin.H{
			"message":"Welcome to LedgerX",
		})
	})

	fmt.Printf("server running on port %s\n",port)

	return router
}