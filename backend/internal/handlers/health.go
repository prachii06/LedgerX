package handlers

import "github.com/gin-gonic/gin"

func Health(c *gin.Context){
	c.JSON(200, gin.H{
		"status": "healthy",
	})
}

func Ready(c *gin.Context){
	c.JSON(200, gin.H{
		"status": "ready",
	})
}

func Live(c *gin.Context){
	c.JSON(200, gin.H{
		"status": "alive",
	})
}