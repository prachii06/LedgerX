package services

import (
	"context"
	"fmt"

	"github.com/prachii06/LedgerX/internal/models"
)

type TransactionGenerator interface {
	Generate(count int) ([]models.GeneratedTransaction, error)
}

type EventGenerator interface {
	GenerateForTransaction(
		ctx context.Context,
		transactionID string,
		amount float64,
		currency string,
	) error
}

type SimulationService struct {
	transactionGenerator TransactionGenerator
	eventGenerator       EventGenerator
}

func NewSimulationService(
	transactionGenerator TransactionGenerator,
	eventGenerator EventGenerator,
) *SimulationService {

	return &SimulationService{
		transactionGenerator: transactionGenerator,
		eventGenerator:       eventGenerator,
	}
}

func (s *SimulationService) Generate(count int) ([]string, error) {

	transactions, err := s.transactionGenerator.Generate(count)
	if err != nil {
		return nil, err
	}

	fmt.Println("TRANSACTIONS:", transactions)

	var ids []string
	for _, transaction := range transactions {
		ids = append(ids, transaction.ID)

		fmt.Println("GENERATING EVENTS FOR:", transaction.ID)

		err := s.eventGenerator.GenerateForTransaction(
			context.Background(),
			transaction.ID,
			transaction.Amount,
			transaction.Currency,
		)

		if err != nil {
			return ids, err
		}
	}

	return ids, nil
}
