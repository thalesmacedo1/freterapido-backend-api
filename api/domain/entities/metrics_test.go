package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	domain "github.com/thalesmacedo1/freterapido-backend-api/api/domain/entities"
)

func TestQuoteMetrics(t *testing.T) {
	metrics := domain.QuoteMetrics{
		CarrierName:          "EXPRESSO FR",
		TotalQuotes:          10,
		TotalShippingPrice:   150.50,
		AverageShippingPrice: 15.05,
	}

	assert.Equal(t, "EXPRESSO FR", metrics.CarrierName)
	assert.Equal(t, 10, metrics.TotalQuotes)
	assert.Equal(t, 150.50, metrics.TotalShippingPrice)
	assert.Equal(t, 15.05, metrics.AverageShippingPrice)
}

func TestMetricsResponse(t *testing.T) {
	metrics1 := domain.QuoteMetrics{
		CarrierName:          "EXPRESSO FR",
		TotalQuotes:          10,
		TotalShippingPrice:   150.50,
		AverageShippingPrice: 15.05,
	}

	metrics2 := domain.QuoteMetrics{
		CarrierName:          "Correios",
		TotalQuotes:          5,
		TotalShippingPrice:   120.25,
		AverageShippingPrice: 24.05,
	}

	cheapestAndMostExpensive := domain.CheapestAndMostExpensive{
		CheapestShipping:      12.50,
		MostExpensiveShipping: 30.75,
	}

	response := domain.MetricsResponse{
		CarrierMetrics:           []domain.QuoteMetrics{metrics1, metrics2},
		CheapestAndMostExpensive: cheapestAndMostExpensive,
	}

	assert.Len(t, response.CarrierMetrics, 2)
	assert.Equal(t, "EXPRESSO FR", response.CarrierMetrics[0].CarrierName)
	assert.Equal(t, 10, response.CarrierMetrics[0].TotalQuotes)
	assert.Equal(t, 150.50, response.CarrierMetrics[0].TotalShippingPrice)
	assert.Equal(t, 15.05, response.CarrierMetrics[0].AverageShippingPrice)

	assert.Equal(t, "Correios", response.CarrierMetrics[1].CarrierName)
	assert.Equal(t, 5, response.CarrierMetrics[1].TotalQuotes)
	assert.Equal(t, 120.25, response.CarrierMetrics[1].TotalShippingPrice)
	assert.Equal(t, 24.05, response.CarrierMetrics[1].AverageShippingPrice)

	assert.Equal(t, 12.50, response.CheapestAndMostExpensive.CheapestShipping)
	assert.Equal(t, 30.75, response.CheapestAndMostExpensive.MostExpensiveShipping)
}
