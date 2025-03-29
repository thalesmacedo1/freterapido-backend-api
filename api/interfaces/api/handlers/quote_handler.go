package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thalesmacedo1/freterapido-backend-api/application/usecases"
	"github.com/thalesmacedo1/freterapido-backend-api/infrastructure/logger"
)

type VaccineHandler struct {
	GetMetricsUC usecases.GetMetricsUseCase
	Logger       logger.Logger
}

func NewQuoteHandler(GetMetricsUC usecases.GetMetricsUseCase, logger logger.Logger) *VaccineHandler {
	return &VaccineHandler{
		GetMetricsUC: GetMetricsUC,
		Logger:       logger,
	}
}

// GetMetrics godoc
// @Summary Retrieve vaccines used in a country
// @Description Fetches the list of vaccines used in a specific country
// @Tags Vaccine
// @Accept json
// @Produce json
// @Param countryCode path string true "ISO Country Code"
// @Success 200 {object} usecases.GetMetricsOutput "Successful retrieval of vaccines used"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/countries/{countryCode}/vaccines [get]
func (h *VaccineHandler) MakeQuotation(c *gin.Context) {
	countryCode := c.Param("countryCode")

	input := usecases.GetMetricsInput{
		CountryCode: countryCode,
	}

	output, err := h.GetMetricsUC.Execute(c.Request.Context(), input)
	if err != nil {
		h.Logger.Errorf("Error executing GetMetricsUseCase: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve vaccines used."})
		return
	}

	c.JSON(http.StatusOK, output)
}
