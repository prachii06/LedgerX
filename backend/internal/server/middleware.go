package server

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prachii06/LedgerX/internal/metrics"
)

// MetricsMiddleware tracks HTTP request latency, counts, and status codes.
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())
		method := c.Request.Method
		route := c.FullPath()

		if route == "" {
			route = "unknown"
		}

		// Skip metrics endpoint itself to avoid noise
		if route == "/metrics" {
			return
		}

		metrics.HTTPRequestsTotal.WithLabelValues(method, route, status).Inc()
		metrics.HTTPRequestDuration.WithLabelValues(method, route).Observe(duration)
	}
}
