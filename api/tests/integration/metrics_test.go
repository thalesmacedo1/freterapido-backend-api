package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	domain "github.com/thalesmacedo1/freterapido-backend-api/api/domain/entities"
	"gorm.io/gorm"
)

func TestMetricsEndpoint_Integration(t *testing.T) {
	// Skip if test environment is not set up
	if testRouter == nil {
		t.Skip("Test environment not set up")
	}

	// Create test data in database
	now := time.Now()
	quote1 := &domain.QuoteResponse{
		Model: gorm.Model{
			CreatedAt: now,
			UpdatedAt: now,
		},
		Carriers: []domain.Carrier{
			{
				Name:     "EXPRESSO FR",
				Service:  "Rodoviário",
				Deadline: "3",
				Price:    17.0,
			},
		},
	}

	quote2 := &domain.QuoteResponse{
		Model: gorm.Model{
			CreatedAt: now.Add(-1 * time.Hour),
			UpdatedAt: now.Add(-1 * time.Hour),
		},
		Carriers: []domain.Carrier{
			{
				Name:     "Correios",
				Service:  "SEDEX",
				Deadline: "1",
				Price:    25.0,
			},
		},
	}

	// Save quotes to database
	err := testQuoteRepository.SaveQuote(context.Background(), quote1)
	assert.NoError(t, err)

	err = testQuoteRepository.SaveQuote(context.Background(), quote2)
	assert.NoError(t, err)

	// Create HTTP request
	req, err := http.NewRequest("GET", "/metrics?last_quotes=10", nil)
	assert.NoError(t, err)

	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	testRouter.ServeHTTP(w, req)

	// Check response status
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse response body
	var response domain.MetricsResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Check response content
	assert.Len(t, response.CarrierMetrics, 2) // Should have 2 carriers
	assert.Equal(t, 17.0, response.CheapestAndMostExpensive.CheapestShipping)
	assert.Equal(t, 25.0, response.CheapestAndMostExpensive.MostExpensiveShipping)

	// Validate carriers metrics
	for _, carrier := range response.CarrierMetrics {
		if carrier.CarrierName == "EXPRESSO FR" {
			assert.Equal(t, 1, carrier.TotalQuotes)
			assert.Equal(t, 17.0, carrier.TotalShippingPrice)
			assert.Equal(t, 17.0, carrier.AverageShippingPrice)
		} else if carrier.CarrierName == "Correios" {
			assert.Equal(t, 1, carrier.TotalQuotes)
			assert.Equal(t, 25.0, carrier.TotalShippingPrice)
			assert.Equal(t, 25.0, carrier.AverageShippingPrice)
		}
	}
}
