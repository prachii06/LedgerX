package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prachii06/LedgerX/internal/handlers"
	"github.com/prachii06/LedgerX/internal/redis"
	"github.com/prachii06/LedgerX/internal/websocket"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func New(port string, corsAllowedOrigins string, db *pgxpool.Pool, redisClient *redis.Client, kafkaBrokers string, transactionHandler *handlers.TransactionHandler, simulationHandler *handlers.SimulationHandler, reconciliationHandler *handlers.ReconciliationHandler, eventHandler *handlers.EventHandler, dashboardHandler *handlers.DashboardHandler, hub *websocket.Hub) *gin.Engine {
	router := gin.Default()

	// Add Prometheus HTTP middleware
	router.Use(MetricsMiddleware())

	// Add CORS middleware
	corsMiddleware := func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		allowedOrigin := "http://localhost:3000" // Default for local dev
		if corsAllowedOrigins != "" {
			allowedOrigin = corsAllowedOrigins
		}

		// Always allow localhost in development, or match the exact origin if specified
		if origin == "http://localhost:3000" || origin == allowedOrigin {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		} else if allowedOrigin == "*" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		} else {
			// Fallback (might fail CORS if origin doesn't match)
			c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		}
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}

	router.Use(corsMiddleware)

	// Ensure OPTIONS requests don't instantly 404/405 before hitting the middleware
	router.NoRoute(corsMiddleware)
	router.NoMethod(corsMiddleware)

	healthHandler := handlers.NewHealthHandler(db, redisClient, kafkaBrokers)

	router.GET("/health", healthHandler.Health)
	router.GET("/ready", healthHandler.Ready)
	router.GET("/live", healthHandler.Live)
	router.GET("/transactions", transactionHandler.List)
	router.POST("/transactions", transactionHandler.Create)
	router.GET("/transactions/:id", transactionHandler.Get)
	router.GET("/events", eventHandler.List)
	router.GET("/dashboard/overview", dashboardHandler.Overview)
	router.POST("/simulate", simulationHandler.Generate)
	router.GET("/reconcile/:transaction_id", reconciliationHandler.Reconcile)
	router.GET("/reconcile/:transaction_id/history", reconciliationHandler.GetHistory)
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Register WebSocket endpoint
	router.GET("/ws", func(c *gin.Context) {
		hub.ServeWs(c)
	})

	fmt.Printf("server running on port %s\n", port)

	return router
}
