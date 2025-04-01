package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thalesmacedo1/freterapido-backend-api/api/application/usecases"
	domain "github.com/thalesmacedo1/freterapido-backend-api/api/domain/entities"
)

type QuoteController struct {
	getShippingQuotationUseCase *usecases.GetShippingQuotationUseCase
}

func NewQuoteController(getShippingQuotationUseCase *usecases.GetShippingQuotationUseCase) *QuoteController {
	return &QuoteController{
		getShippingQuotationUseCase: getShippingQuotationUseCase,
	}
}

// GetQuote obtém cotações de frete de diferentes transportadoras
// @Summary Obter cotações de frete
// @Description Retorna cotações de frete de diferentes transportadoras com base nos dados enviados
// @Tags cotações
// @Accept json
// @Produce json
// @Param request body domain.QuoteRequest true "Dados para cotação de frete"
// @Success 200 {object} domain.QuoteResponse "Cotações de frete disponíveis"
// @Failure 400 {object} map[string]string "Erro de requisição inválida"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Router /quote [post]
func (c *QuoteController) GetQuote(ctx *gin.Context) {
	var request domain.QuoteRequest

	// Bind JSON request
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format: " + err.Error()})
		return
	}

	// Validate request
	if err := validateQuoteRequest(request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Execute use case
	response, err := c.getShippingQuotationUseCase.Execute(ctx, request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get shipping quotation: " + err.Error()})
		return
	}

	// Return response
	ctx.JSON(http.StatusOK, response)
}

func validateQuoteRequest(request domain.QuoteRequest) error {
	// Basic validation logic can be expanded
	if request.Recipient.Address.Zipcode == "" {
		return &InputError{Field: "recipient.address.zipcode", Message: "Zipcode cannot be empty"}
	}

	if len(request.Volumes) == 0 {
		return &InputError{Field: "volumes", Message: "At least one volume is required"}
	}

	return nil
}

type InputError struct {
	Field   string
	Message string
}

func (e *InputError) Error() string {
	return e.Message + " for field " + e.Field
}
