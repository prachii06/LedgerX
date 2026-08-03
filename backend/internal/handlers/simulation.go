package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/prachii06/LedgerX/internal/generator"
	"github.com/prachii06/LedgerX/internal/models"
)

type SimulationHandler struct {
	generator *generator.TransactionGenerator
}

func NewSimulationHandler(
	generator *generator.TransactionGenerator,
) *SimulationHandler {

	return &SimulationHandler{
		generator: generator,
	}
}

func (h *SimulationHandler) Generate(c *gin.Context) {

	var request models.SimulationRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := h.generator.Generate(request.Count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Simulation completed",
		"count":   request.Count,
	})
}