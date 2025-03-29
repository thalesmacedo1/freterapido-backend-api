package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thalesmacedo1/covid-api/application/usecases"
	"github.com/thalesmacedo1/covid-api/infrastructure/logger"
)

type MetricsHandler struct {
	GetVaccinesUsedUC usecases.GetVaccinesUsedUseCase
	Logger            logger.Logger
}

func NewMetricsHandler(getVaccinesUsedUC usecases.GetMetricsUC, logger logger.Logger) *MetricsHandler {
	return &MetricsHandler{
		GetVaccinesUsedUC: getVaccinesUsedUC,
		Logger:            logger,
	}
}

// GetVaccinesUsed godoc
// @Summary Retrieve vaccines used in a country
// @Description Fetches the list of vaccines used in a specific country
// @Tags Vaccine
// @Accept json
// @Produce json
// @Param countryCode path string true "ISO Country Code"
// @Success 200 {object} usecases.GetVaccinesUsedOutput "Successful retrieval of vaccines used"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/countries/{countryCode}/vaccines [get]
func (h *VaccineHandler) GetVaccinesUsed(c *gin.Context) {
	countryCode := c.Param("countryCode")

	input := usecases.GetVaccinesUsedInput{
		CountryCode: countryCode,
	}

	output, err := h.GetVaccinesUsedUC.Execute(c.Request.Context(), input)
	if err != nil {
		h.Logger.Errorf("Error executing GetVaccinesUsedUseCase: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve vaccines used."})
		return
	}

	c.JSON(http.StatusOK, output)
}
