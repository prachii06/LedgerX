package generator

import (
	"context"
	"fmt"
	"math/rand"
	"github.com/prachii06/LedgerX/internal/models"
	"github.com/prachii06/LedgerX/internal/services"
)

type TransactionGenerator struct {
	service *services.TransactionService
}

func NewTransactionGenerator(service *services.TransactionService) *TransactionGenerator {
	return &TransactionGenerator{
		service: service,
	}
}

func (g *TransactionGenerator) Generate(count int) error {

	for i := 1; i <= count; i++ {

		transaction := &models.Transaction{
			ExternalID: fmt.Sprintf("SIM-%06d", i),
			Amount:     float64(rand.Intn(9900) + 100),
			Currency:   "INR",
		}

		err := g.service.CreateTransaction(context.Background(), transaction)
		if err != nil {
			return err
		}
	}

	return nil
}