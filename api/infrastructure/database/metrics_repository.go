package database

import (
	"context"
	"math"

	domain "github.com/thalesmacedo1/freterapido-backend-api/api/domain/entities"
	"gorm.io/gorm"
)

type MetricsRepositoryImpl struct {
	db *gorm.DB
}

func NewMetricsRepository(db *gorm.DB) domain.MetricsRepository {
	return &MetricsRepositoryImpl{
		db: db,
	}
}

func (r *MetricsRepositoryImpl) GetMetrics(ctx context.Context, lastQuotes int) (*domain.MetricsResponse, error) {
	// Get the quotes from the database
	var quoteRepo = &QuoteRepositoryImpl{db: r.db}
	quotes, err := quoteRepo.GetLastQuotes(ctx, lastQuotes)
	if err != nil {
		return nil, err
	}

	// Initialize response
	response := &domain.MetricsResponse{
		CarrierMetrics: []domain.QuoteMetrics{},
		CheapestAndMostExpensive: domain.CheapestAndMostExpensive{
			CheapestShipping:      math.MaxFloat64,
			MostExpensiveShipping: 0,
		},
	}

	// Map to store metrics per carrier
	carrierMetrics := make(map[string]*domain.QuoteMetrics)

	// Process each quote
	for _, quote := range quotes {
		for _, carrier := range quote.Carriers {
			// Update cheapest and most expensive shipping
			if carrier.Price < response.CheapestAndMostExpensive.CheapestShipping {
				response.CheapestAndMostExpensive.CheapestShipping = carrier.Price
			}
			if carrier.Price > response.CheapestAndMostExpensive.MostExpensiveShipping {
				response.CheapestAndMostExpensive.MostExpensiveShipping = carrier.Price
			}

			// Update carrier metrics
			if _, exists := carrierMetrics[carrier.Name]; !exists {
				carrierMetrics[carrier.Name] = &domain.QuoteMetrics{
					CarrierName:          carrier.Name,
					TotalQuotes:          0,
					TotalShippingPrice:   0,
					AverageShippingPrice: 0,
				}
			}

			carrierMetrics[carrier.Name].TotalQuotes++
			carrierMetrics[carrier.Name].TotalShippingPrice += carrier.Price
		}
	}

	// Calculate averages and build response
	for _, metrics := range carrierMetrics {
		if metrics.TotalQuotes > 0 {
			metrics.AverageShippingPrice = metrics.TotalShippingPrice / float64(metrics.TotalQuotes)
		}
		response.CarrierMetrics = append(response.CarrierMetrics, *metrics)
	}

	// If no quotes found, set cheapest to 0
	if response.CheapestAndMostExpensive.CheapestShipping == math.MaxFloat64 {
		response.CheapestAndMostExpensive.CheapestShipping = 0
	}

	return response, nil
}
