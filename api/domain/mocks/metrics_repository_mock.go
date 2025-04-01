package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	domain "github.com/thalesmacedo1/freterapido-backend-api/api/domain/entities"
)

// MockMetricsRepository is a mock implementation of the MetricsRepository interface
type MockMetricsRepository struct {
	mock.Mock
}

// GetMetrics is a mock implementation of the GetMetrics method
func (m *MockMetricsRepository) GetMetrics(ctx context.Context, lastQuotes int) (*domain.MetricsResponse, error) {
	args := m.Called(ctx, lastQuotes)

	// If the return value is nil, return nil to avoid casting nil to *domain.MetricsResponse
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.MetricsResponse), args.Error(1)
}
