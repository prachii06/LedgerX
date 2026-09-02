package handlers

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HealthHandler struct {
	DB *pgxpool.Pool
}

func NewHealthHandler(db *pgxpool.Pool) *HealthHandler {
	return &HealthHandler{
		DB: db,
	}
}

func (h *HealthHandler) Live(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "alive",
	})
}

func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "healthy",
	})
}

func (h *HealthHandler) Ready(c *gin.Context) {

	err := h.DB.Ping(context.Background())

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
