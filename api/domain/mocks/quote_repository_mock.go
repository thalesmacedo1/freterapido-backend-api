package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	domain "github.com/thalesmacedo1/freterapido-backend-api/api/domain/entities"
)

// MockQuoteRepository is a mock implementation of the QuoteRepository interface
type MockQuoteRepository struct {
	mock.Mock
}

// SaveQuote is a mock implementation of the SaveQuote method
func (m *MockQuoteRepository) SaveQuote(ctx context.Context, quote *domain.QuoteResponse) error {
	args := m.Called(ctx, quote)
	return args.Error(0)
}

// GetLastQuotes is a mock implementation of the GetLastQuotes method
func (m *MockQuoteRepository) GetLastQuotes(ctx context.Context, limit int) ([]domain.QuoteResponse, error) {
	args := m.Called(ctx, limit)

	// If the return value is nil, return nil to avoid casting nil to []domain.QuoteResponse
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]domain.QuoteResponse), args.Error(1)
}
