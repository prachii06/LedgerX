package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/prachii06/LedgerX/internal/models"
	"github.com/prachii06/LedgerX/internal/services"
	"net/http"
)

type SimulationHandler struct {
	service *services.SimulationService
}

func NewSimulationHandler(
	service *services.SimulationService,
) *SimulationHandler {

	return &SimulationHandler{
		service: service,
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

	ids, err := h.service.Generate(request.Count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":         "Simulation completed",
		"count":           request.Count,
		"transaction_ids": ids,
	})
}
