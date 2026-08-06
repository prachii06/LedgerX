package services

type Generator interface {
	Generate(count int) error
}

type SimulationService struct {
	generator Generator
}

func NewSimulationService(
	generator Generator,
) *SimulationService {

	return &SimulationService{
		generator: generator,
	}
}

func (s *SimulationService) Generate(count int) error {
	return s.generator.Generate(count)
}
