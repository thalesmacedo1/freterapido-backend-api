package database

import (
	"context"
	"fmt"
	"math"
	"sort"
	"sync"

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

	// Track metrics per carrier with a mutex for thread safety
	carrierMetrics := make(map[string]*domain.QuoteMetrics)
	var metricsMutex sync.Mutex
	var cheapestMutex sync.Mutex

	// Use a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Process quotes concurrently
	for i, quote := range quotes {
		// Skip quotes with no carriers
		if len(quote.Carriers) == 0 {
			fmt.Printf("Warning: Quote ID %d has no carriers\n", quote.ID)
			continue
		}

		// Add to wait group
		wg.Add(1)

		// Process each quote in a separate goroutine
		go func(index int, q domain.QuoteResponse) {
			defer wg.Done()

			// Check for context cancellation
			select {
			case <-ctx.Done():
				// Context was cancelled, stop processing
				return
			default:
				// Continue processing
			}

			// Process all carriers in this quote
			for j, carrier := range q.Carriers {
				// Check for context cancellation periodically
				select {
				case <-ctx.Done():
					return
				default:
					// Continue processing
				}

				fmt.Printf("Processing quote %d, carrier %d: %s with price %.2f\n",
					index, j, carrier.Name, carrier.Price)

				// Update cheapest/most expensive with mutex protection
				cheapestMutex.Lock()
				if carrier.Price < response.CheapestAndMostExpensive.CheapestShipping {
					response.CheapestAndMostExpensive.CheapestShipping = carrier.Price
				}
				if carrier.Price > response.CheapestAndMostExpensive.MostExpensiveShipping {
					response.CheapestAndMostExpensive.MostExpensiveShipping = carrier.Price
				}
				cheapestMutex.Unlock()

				// Update carrier metrics with mutex protection
				metricsMutex.Lock()
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
				metricsMutex.Unlock()
			}
		}(i, quote)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	// Convert map to slice in a deterministic order (alphabetical by carrier name)
	var carrierNames []string
	for name := range carrierMetrics {
		carrierNames = append(carrierNames, name)
	}
	sort.Strings(carrierNames)

	// Use channels to collect metrics in parallel
	metricsChan := make(chan domain.QuoteMetrics, len(carrierNames))
	var avgWg sync.WaitGroup

	// Calculate average prices in parallel
	for _, name := range carrierNames {
		avgWg.Add(1)
		go func(name string, metrics *domain.QuoteMetrics) {
			defer avgWg.Done()

			// Clone the metrics to avoid concurrent modification
			metricsCopy := *metrics

			if metricsCopy.TotalQuotes > 0 {
				metricsCopy.AverageShippingPrice = metricsCopy.TotalShippingPrice / float64(metricsCopy.TotalQuotes)
			}

			// Send to channel
			metricsChan <- metricsCopy
		}(name, carrierMetrics[name])
	}

	// Close channel when all goroutines are done
	go func() {
		avgWg.Wait()
		close(metricsChan)
	}()

	// Collect results from channel
	for metric := range metricsChan {
		response.CarrierMetrics = append(response.CarrierMetrics, metric)
	}

	// Sort the final metrics slice to maintain deterministic order
	sort.Slice(response.CarrierMetrics, func(i, j int) bool {
		return response.CarrierMetrics[i].CarrierName < response.CarrierMetrics[j].CarrierName
	})

	// Handle edge case when no valid carriers were found
	if response.CheapestAndMostExpensive.CheapestShipping == math.MaxFloat64 {
		response.CheapestAndMostExpensive.CheapestShipping = 0
	}

	return response, nil
}
