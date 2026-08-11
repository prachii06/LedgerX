package generator

import (
	"context"
	"fmt"
	"math/rand"
	"time"
	"github.com/prachii06/LedgerX/internal/models"
)

type TransactionCreator interface {
	CreateTransaction(
		ctx context.Context,
		transaction *models.Transaction,
	) error
}

type TransactionGenerator struct {
	service TransactionCreator
}

func NewTransactionGenerator(
	service TransactionCreator,
) *TransactionGenerator {
	return &TransactionGenerator{
		service: service,
	}
}

func (g *TransactionGenerator) Generate(count int) ([]models.GeneratedTransaction, error) {

	transactions := make([]models.GeneratedTransaction, 0, count)

	for i := 1; i <= count; i++ {

		transaction := &models.Transaction{
			ExternalID: fmt.Sprintf("SIM-%d-%06d", time.Now().UnixNano(), i),
			Amount:     float64(rand.Intn(9900) + 100),
			Currency:   "INR",
		}

		err := g.service.CreateTransaction(
			context.Background(),
			transaction,
		)

		if err != nil {
			return nil, err
		}

		transactions = append(
			transactions,
			models.GeneratedTransaction{
				ID:       transaction.ID,
				Amount:   transaction.Amount,
				Currency: transaction.Currency,
			},
		)
	}

	return transactions, nil
}