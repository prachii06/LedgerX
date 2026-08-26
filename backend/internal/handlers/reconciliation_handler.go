package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/prachii06/LedgerX/internal/services"
)

type ReconciliationHandler struct {
	service *services.ReconciliationService
}

func NewReconciliationHandler(
	service *services.ReconciliationService,
) *ReconciliationHandler {
	return &ReconciliationHandler{
		service: service,
	}
}

func (h *ReconciliationHandler) Reconcile(c *gin.Context) {

	transactionID := c.Param("transaction_id")

	result, err := h.service.ReconcileTransaction(
		c.Request.Context(),
		transactionID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}


func (h *ReconciliationHandler) GetHistory(c *gin.Context) {

	transactionID := c.Param("transaction_id")

	history, err := h.service.GetReconciliationHistory(
		c.Request.Context(),
		transactionID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, history)
}