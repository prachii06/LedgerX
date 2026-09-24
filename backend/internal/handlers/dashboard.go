package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/prachii06/LedgerX/internal/services"
	"net/http"
)

type DashboardHandler struct {
	service *services.TransactionService
}

func NewDashboardHandler(service *services.TransactionService) *DashboardHandler {
	return &DashboardHandler{
		service: service,
	}
}

func (h *DashboardHandler) Overview(c *gin.Context) {
	total, reconciled, issues, pending, err := h.service.GetOverview(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total_transactions": total,
		"reconciled":         reconciled,
		"issues":             issues,
		"pending":            pending,
	})
}
