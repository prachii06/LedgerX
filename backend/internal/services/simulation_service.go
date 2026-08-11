package services

import(

"context"
"fmt"

)

type TransactionGenerator interface {
	Generate(count int) ([]string, error)
}

type EventGenerator interface {
	GenerateForTransaction(
		ctx context.Context,
		transactionID string,
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

func (s *SimulationService) Generate(count int) error {

	transactionIDs, err := s.transactionGenerator.Generate(count)
	if err != nil {
		return err
	}

	fmt.Println("TRANSACTION IDS:",transactionIDs)

	for _, transactionID := range transactionIDs {

		fmt.Println("GENERATING EVENTS FOR:", transactionID)

		err := s.eventGenerator.GenerateForTransaction(
			context.Background(),
			transactionID,
		)

		if err != nil {
			return err
		}
	}

	return nil
}