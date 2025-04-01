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
	// Get metrics from repository
	return uc.metricsRepository.GetMetrics(ctx, lastQuotes)
}
