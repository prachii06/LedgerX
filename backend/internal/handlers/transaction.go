package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/prachii06/LedgerX/internal/models"
	"github.com/prachii06/LedgerX/internal/services"
	"net/http"
	"strconv"
)

type TransactionHandler struct {
	service *services.TransactionService
}

func NewTransactionHandler(service *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		service: service,
	}
}

func (h *TransactionHandler) Create(c *gin.Context) {
	var transaction models.Transaction

	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := h.service.CreateTransaction(c.Request.Context(), &transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, transaction)
}

func (h *TransactionHandler) List(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	transactions, err := h.service.GetTransactions(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if transactions == nil {
		transactions = []models.Transaction{}
	}

	c.JSON(http.StatusOK, transactions)
}

func (h *TransactionHandler) Get(c *gin.Context) {
	id := c.Param("id")
	transaction, err := h.service.GetTransaction(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	c.JSON(http.StatusOK, transaction)
}
