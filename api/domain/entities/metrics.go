package domain

import (
	"context"
)

// QuoteMetrics representa métricas para uma transportadora específica
// @Description Métricas de cotações para uma transportadora específica
type QuoteMetrics struct {
	// Nome da transportadora
	// @example "EXPRESSO FR"
	CarrierName string `json:"carrier_name"`
	// Total de cotações realizadas
	// @example 10
	TotalQuotes int `json:"total_quotes"`
	// Valor total dos fretes cotados
	// @example 150.50
	TotalShippingPrice float64 `json:"total_shipping_price"`
	// Valor médio dos fretes cotados
	// @example 15.05
	AverageShippingPrice float64 `json:"average_shipping_price"`
}

// MetricsResponse é a resposta completa de métricas
// @Description Resposta completa com todas as métricas de cotações
type MetricsResponse struct {
	// Métricas por transportadora
	// @Description Lista de métricas por transportadora
	CarrierMetrics []QuoteMetrics `json:"carrier_metrics"`
	// Informações sobre cotações mais baratas e mais caras
	// @Description Detalhes sobre os valores mínimos e máximos de frete
	CheapestAndMostExpensive CheapestAndMostExpensive `json:"cheapest_and_most_expensive"`
}

// CheapestAndMostExpensive representa os fretes mais baratos e mais caros
// @Description Valores mínimos e máximos encontrados nas cotações
type CheapestAndMostExpensive struct {
	// Valor do frete mais barato
	// @example 12.50
	CheapestShipping float64 `json:"cheapest_shipping"`
	// Valor do frete mais caro
	// @example 30.75
	MostExpensiveShipping float64 `json:"most_expensive_shipping"`
}

// MetricsRepository define a interface para operações de métricas
type MetricsRepository interface {
	GetMetrics(ctx context.Context, lastQuotes int) (*MetricsResponse, error)
}
