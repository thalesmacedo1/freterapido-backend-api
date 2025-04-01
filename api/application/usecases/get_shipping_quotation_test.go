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

// This is a simplified test focusing on the repository interaction
// A more comprehensive test would also mock the HTTP client for testing the API call
func TestGetShippingQuotationUseCase_Execute(t *testing.T) {
	// Skip this test since we need to mock HTTP client functionality
	// which would require refactoring the use case to accept a client interface
	t.Skip("Skipping test that requires HTTP client mocking")

	// Create a mock repository
	mockRepo := new(mocks.MockQuoteRepository)

	// Setup test request
	request := domain.QuoteRequest{}
	request.Recipient.Address.Zipcode = "01311000"

	volume := domain.Volume{
		Category:      7,
		Amount:        1,
		UnitaryWeight: 5.0,
		Price:         349.0,
		SKU:           "abc-teste-123",
		Height:        0.2,
		Width:         0.2,
		Length:        0.2,
	}

	request.Volumes = append(request.Volumes, volume)

	// Setup mock for repository save
	mockRepo.On("SaveQuote", mock.Anything, mock.AnythingOfType("*domain.QuoteResponse")).Return(nil)

	// Create the use case with the mock repository
	useCase := usecases.NewGetShippingQuotationUseCase(mockRepo)

	// Execute the use case
	result, err := useCase.Execute(context.Background(), request)

	// Assert results
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Verify expectations were met
	mockRepo.AssertExpectations(t)
}

// Test error handling when saving quote
func TestGetShippingQuotationUseCase_SaveError(t *testing.T) {
	// Skip this test since we need to mock HTTP client functionality
	t.Skip("Skipping test that requires HTTP client mocking")

	// Create a mock repository
	mockRepo := new(mocks.MockQuoteRepository)

	// Setup test request
	request := domain.QuoteRequest{}
	request.Recipient.Address.Zipcode = "01311000"

	volume := domain.Volume{
		Category:      7,
		Amount:        1,
		UnitaryWeight: 5.0,
		Price:         349.0,
		SKU:           "abc-teste-123",
		Height:        0.2,
		Width:         0.2,
		Length:        0.2,
	}

	request.Volumes = append(request.Volumes, volume)

	// Setup mock for repository save with error
	expectedError := errors.New("database error")
	mockRepo.On("SaveQuote", mock.Anything, mock.AnythingOfType("*domain.QuoteResponse")).Return(expectedError)

	// Create the use case with the mock repository
	useCase := usecases.NewGetShippingQuotationUseCase(mockRepo)

	// Execute the use case
	result, err := useCase.Execute(context.Background(), request)

	// Assert results
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "error saving quote")

	// Verify expectations were met
	mockRepo.AssertExpectations(t)
}
