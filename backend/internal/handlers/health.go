package handlers

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Health(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "healthy",
	})
}

func Live(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "alive",
	})
}

func Ready(db *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {

		err := db.Ping(context.Background())

		if err != nil {
			c.JSON(503, gin.H{
				"status": "database unavailable",
			})
			return
		}

		c.JSON(200, gin.H{
			"status": "ready",
		})
	}
}