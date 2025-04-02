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
	var quoteRepo = &QuoteRepositoryImpl{db: r.db}
	quotes, err := quoteRepo.GetLastQuotes(ctx, lastQuotes)
	if err != nil {
		return nil, err
	}

	response := &domain.MetricsResponse{
		CarrierMetrics: []domain.QuoteMetrics{},
		CheapestAndMostExpensive: domain.CheapestAndMostExpensive{
			CheapestShipping:      math.MaxFloat64,
			MostExpensiveShipping: 0,
		},
	}

	carrierMetrics := make(map[string]*domain.QuoteMetrics)

	for _, quote := range quotes {
		for _, carrier := range quote.Carriers {
			if carrier.Price < response.CheapestAndMostExpensive.CheapestShipping {
				response.CheapestAndMostExpensive.CheapestShipping = carrier.Price
			}
			if carrier.Price > response.CheapestAndMostExpensive.MostExpensiveShipping {
				response.CheapestAndMostExpensive.MostExpensiveShipping = carrier.Price
			}

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

	for _, metrics := range carrierMetrics {
		if metrics.TotalQuotes > 0 {
			metrics.AverageShippingPrice = metrics.TotalShippingPrice / float64(metrics.TotalQuotes)
		}
		response.CarrierMetrics = append(response.CarrierMetrics, *metrics)
	}

	if response.CheapestAndMostExpensive.CheapestShipping == math.MaxFloat64 {
		response.CheapestAndMostExpensive.CheapestShipping = 0
	}

	return response, nil
}
