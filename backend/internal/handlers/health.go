package handlers

import (
	"context"
	"net"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prachii06/LedgerX/internal/redis"
)

type HealthHandler struct {
	DB           *pgxpool.Pool
	RedisClient  *redis.Client
	KafkaBrokers string
}

func NewHealthHandler(db *pgxpool.Pool, redisClient *redis.Client, kafkaBrokers string) *HealthHandler {
	return &HealthHandler{
		DB:           db,
		RedisClient:  redisClient,
		KafkaBrokers: kafkaBrokers,
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
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	pgStatus := "ok"
	if err := h.DB.Ping(ctx); err != nil {
		pgStatus = "unavailable"
	}

	redisStatus := "ok"
	if h.RedisClient != nil {
		if err := h.RedisClient.Ping(ctx); err != nil {
			redisStatus = "unavailable"
		}
	} else {
		redisStatus = "unavailable"
	}

	kafkaStatus := "ok"
	if h.KafkaBrokers != "" {
		conn, err := net.DialTimeout("tcp", h.KafkaBrokers, 2*time.Second)
		if err != nil {
			kafkaStatus = "unavailable"
		} else {
			conn.Close()
		}
	} else {
		kafkaStatus = "unavailable"
	}

	status := "ok"
	if pgStatus != "ok" || redisStatus != "ok" || kafkaStatus != "ok" {
		status = "degraded"
	}

	statusCode := 200
	if status == "degraded" {
		statusCode = 503
	}

	c.JSON(statusCode, gin.H{
		"status": status,
		"services": gin.H{
			"backend":    "ok",
			"postgresql": pgStatus,
			"redis":      redisStatus,
			"kafka":      kafkaStatus,
		},
	})
}
