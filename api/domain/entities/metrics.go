package domain

import (
	"context"
	"time"

	domain "github.com/thalesmacedo1/freterapido-backend-api/api/domain/entities"
)

type QuoteMetrics struct {
	CarrierName          string  `json:"carrier_name"`
	TotalQuotes          int     `json:"total_quotes"`
	TotalShippingPrice   float64 `json:"total_shipping_price"`
	AverageShippingPrice float64 `json:"average_shipping_price"`
}

type CheapestAndMostExpensive struct {
	CheapestShipping      float64 `json:"cheapest_shipping"`
	MostExpensiveShipping float64 `json:"most_expensive_shipping"`
}

type MetricsRepository interface {
	GetShippingQuotation(ctx context.Context, countryCode string, date time.Time) (*domain.QuoteMetrics, error)
}
