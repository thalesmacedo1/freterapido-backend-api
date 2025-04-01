package usecases_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thalesmacedo1/freterapido-backend-api/api/application/usecases"
	domain "github.com/thalesmacedo1/freterapido-backend-api/api/domain/entities"
	"github.com/thalesmacedo1/freterapido-backend-api/api/domain/mocks"
)

func TestGetMetricsUseCase_Execute_Success(t *testing.T) {
	// Create a mock repository
	mockRepo := new(mocks.MockMetricsRepository)

	// Setup test data
	mockResponse := &domain.MetricsResponse{
		CarrierMetrics: []domain.QuoteMetrics{
			{
				CarrierName:          "EXPRESSO FR",
				TotalQuotes:          10,
				TotalShippingPrice:   150.50,
				AverageShippingPrice: 15.05,
			},
			{
				CarrierName:          "Correios",
				TotalQuotes:          5,
				TotalShippingPrice:   120.25,
				AverageShippingPrice: 24.05,
			},
		},
		CheapestAndMostExpensive: domain.CheapestAndMostExpensive{
			CheapestShipping:      12.50,
			MostExpensiveShipping: 30.75,
		},
	}

	// Setup expectations
	lastQuotes := 10
	mockRepo.On("GetMetrics", mock.Anything, lastQuotes).Return(mockResponse, nil)

	// Create the use case with the mock repository
	useCase := usecases.NewGetMetricsUseCase(mockRepo)

	// Execute the use case
	result, err := useCase.Execute(context.Background(), lastQuotes)

	// Assert results
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, len(result.CarrierMetrics))
	assert.Equal(t, "EXPRESSO FR", result.CarrierMetrics[0].CarrierName)
	assert.Equal(t, 10, result.CarrierMetrics[0].TotalQuotes)
	assert.Equal(t, 150.50, result.CarrierMetrics[0].TotalShippingPrice)
	assert.Equal(t, 15.05, result.CarrierMetrics[0].AverageShippingPrice)

	assert.Equal(t, "Correios", result.CarrierMetrics[1].CarrierName)
	assert.Equal(t, 5, result.CarrierMetrics[1].TotalQuotes)
	assert.Equal(t, 120.25, result.CarrierMetrics[1].TotalShippingPrice)
	assert.Equal(t, 24.05, result.CarrierMetrics[1].AverageShippingPrice)

	assert.Equal(t, 12.50, result.CheapestAndMostExpensive.CheapestShipping)
	assert.Equal(t, 30.75, result.CheapestAndMostExpensive.MostExpensiveShipping)

	// Verify expectations were met
	mockRepo.AssertExpectations(t)
}

func TestGetMetricsUseCase_Execute_Error(t *testing.T) {
	// Create a mock repository
	mockRepo := new(mocks.MockMetricsRepository)

	// Setup expectations with an error
	lastQuotes := 10
	expectedError := errors.New("database error")
	mockRepo.On("GetMetrics", mock.Anything, lastQuotes).Return(nil, expectedError)

	// Create the use case with the mock repository
	useCase := usecases.NewGetMetricsUseCase(mockRepo)

	// Execute the use case
	result, err := useCase.Execute(context.Background(), lastQuotes)

	// Assert results
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, result)

	// Verify expectations were met
	mockRepo.AssertExpectations(t)
}
