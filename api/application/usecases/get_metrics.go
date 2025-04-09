package usecases

import (
	"context"

	domain "github.com/thalesmacedo1/freterapido-backend-api/api/domain/entities"
)

type GetMetricsUseCase struct {
	metricsRepository domain.MetricsRepository
}

func NewGetMetricsUseCase(metricsRepository domain.MetricsRepository) *GetMetricsUseCase {
	return &GetMetricsUseCase{
		metricsRepository: metricsRepository,
	}
}

func (uc *GetMetricsUseCase) Execute(ctx context.Context, lastQuotes int) (*domain.MetricsResponse, error) {
	// Possivel uso de goroutines
	// Get metrics from repository
	// var wg sync.WaitGroup
	// var metricsResult *domain.MetricsResponse
	// var metricsErr error
	//
	// wg.Add(1)
	// go func() {
	//     defer wg.Done()
	//     metricsResult, metricsErr = uc.metricsRepository.GetMetrics(ctx, lastQuotes)
	// }()
	//
	// // Fetch other data concurrently here
	//
	// wg.Wait()
	// if metricsErr != nil {
	//     return nil, metricsErr
	// }

	return uc.metricsRepository.GetMetrics(ctx, lastQuotes)
}
