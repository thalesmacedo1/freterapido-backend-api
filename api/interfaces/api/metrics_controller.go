package api

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/thalesmacedo1/freterapido-backend-api/api/application/usecases"
	domain "github.com/thalesmacedo1/freterapido-backend-api/api/domain/entities"
)

type MetricsController struct {
	getMetricsUseCase *usecases.GetMetricsUseCase
}

func NewMetricsController(getMetricsUseCase *usecases.GetMetricsUseCase) *MetricsController {
	return &MetricsController{
		getMetricsUseCase: getMetricsUseCase,
	}
}

// GetMetrics retorna métricas sobre as cotações de frete
// @Summary Obter métricas de cotações
// @Description Retorna métricas e estatísticas sobre as cotações de frete realizadas
// @Tags métricas
// @Accept json
// @Produce json
// @Param last_quotes query int false "Número de cotações recentes a considerar (opcional)"
// @Success 200 {object} domain.MetricsResponse "Métricas de cotações"
// @Failure 400 {object} map[string]string "Erro de parâmetro inválido"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Router /metrics [get]
func (c *MetricsController) GetMetrics(ctx *gin.Context) {
	// Parse last_quotes parameter
	lastQuotesStr := ctx.Query("last_quotes")
	lastQuotes := 0

	if lastQuotesStr != "" {
		var err error
		lastQuotes, err = strconv.Atoi(lastQuotesStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "last_quotes must be a valid integer"})
			return
		}

		if lastQuotes < 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "last_quotes must be a positive integer"})
			return
		}

		// Log the last_quotes parameter for debugging
		fmt.Printf("Processing request with last_quotes=%d\n", lastQuotes)
	} else {
		fmt.Println("Processing request without last_quotes parameter")
	}

	// Execute use case
	metrics, err := c.getMetricsUseCase.Execute(ctx, lastQuotes)
	if err != nil {
		fmt.Printf("Error getting metrics: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get metrics: " + err.Error()})
		return
	}

	// Log detailed metrics for debugging
	fmt.Printf("Returning metrics with %d carrier metrics\n", len(metrics.CarrierMetrics))

	// Use a wait group for concurrent logging
	var logWg sync.WaitGroup

	// Log metrics concurrently
	for i, metric := range metrics.CarrierMetrics {
		logWg.Add(1)
		go func(index int, m domain.QuoteMetrics) {
			defer logWg.Done()
			fmt.Printf("  [%d] Carrier: %s, TotalQuotes: %d, TotalPrice: %.2f, AvgPrice: %.2f\n",
				index, m.CarrierName, m.TotalQuotes, m.TotalShippingPrice, m.AverageShippingPrice)
		}(i, metric)
	}

	// Wait for all logging to complete
	logWg.Wait()

	fmt.Printf("Cheapest: %.2f, Most Expensive: %.2f\n",
		metrics.CheapestAndMostExpensive.CheapestShipping,
		metrics.CheapestAndMostExpensive.MostExpensiveShipping)

	// Return response
	ctx.JSON(http.StatusOK, metrics)
}
