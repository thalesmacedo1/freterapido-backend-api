package database

import (
	"context"
	"fmt"
	"math"
	"sort"

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
	var quotes []domain.QuoteResponse

	// Add deterministic ordering with secondary sort on ID to ensure consistent results
	// Also enable debugging to track the queries being executed
	query := r.db.Debug().Order("created_at DESC, id DESC")

	if lastQuotes > 0 {
		query = query.Limit(lastQuotes)
	}

	// Execute the query to find quotes
	if err := query.Find(&quotes).Error; err != nil {
		return nil, err
	}

	// If no quotes found, return empty metrics
	if len(quotes) == 0 {
		return &domain.MetricsResponse{
			CarrierMetrics: []domain.QuoteMetrics{},
			CheapestAndMostExpensive: domain.CheapestAndMostExpensive{
				CheapestShipping:      0,
				MostExpensiveShipping: 0,
			},
		}, nil
	}

	// Debug: print each quote's carriers
	for i, quote := range quotes {
		fmt.Printf("Quote %d has %d carriers\n", i, len(quote.Carriers))
	}

	// Initialize response structure
	response := &domain.MetricsResponse{
		CarrierMetrics: []domain.QuoteMetrics{},
		CheapestAndMostExpensive: domain.CheapestAndMostExpensive{
			CheapestShipping:      math.MaxFloat64,
			MostExpensiveShipping: 0,
		},
	}

	// Track metrics per carrier
	carrierMetrics := make(map[string]*domain.QuoteMetrics)

	// Process quotes and calculate metrics
	for i, quote := range quotes {
		// Skip quotes with no carriers
		if len(quote.Carriers) == 0 {
			fmt.Printf("Warning: Quote ID %d has no carriers\n", quote.ID)
			continue
		}

		for j, carrier := range quote.Carriers {
			fmt.Printf("Processing quote %d, carrier %d: %s with price %.2f\n",
				i, j, carrier.Name, carrier.Price)

			// Update cheapest/most expensive
			if carrier.Price < response.CheapestAndMostExpensive.CheapestShipping {
				response.CheapestAndMostExpensive.CheapestShipping = carrier.Price
			}
			if carrier.Price > response.CheapestAndMostExpensive.MostExpensiveShipping {
				response.CheapestAndMostExpensive.MostExpensiveShipping = carrier.Price
			}

			// Create or update carrier metrics
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

	// Convert map to slice in a deterministic order (alphabetical by carrier name)
	var carrierNames []string
	for name := range carrierMetrics {
		carrierNames = append(carrierNames, name)
	}
	sort.Strings(carrierNames)

	// Build metrics in consistent order
	for _, name := range carrierNames {
		metrics := carrierMetrics[name]
		if metrics.TotalQuotes > 0 {
			metrics.AverageShippingPrice = metrics.TotalShippingPrice / float64(metrics.TotalQuotes)
		}
		response.CarrierMetrics = append(response.CarrierMetrics, *metrics)
	}

	// Handle edge case when no valid carriers were found
	if response.CheapestAndMostExpensive.CheapestShipping == math.MaxFloat64 {
		response.CheapestAndMostExpensive.CheapestShipping = 0
	}

	return response, nil
}
