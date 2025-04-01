package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/thalesmacedo1/freterapido-backend-api/api/application/usecases"
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
	}

	// Execute use case
	metrics, err := c.getMetricsUseCase.Execute(ctx, lastQuotes)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get metrics: " + err.Error()})
		return
	}

	// Return response
	ctx.JSON(http.StatusOK, metrics)
}
