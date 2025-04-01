package api_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	domain "github.com/thalesmacedo1/freterapido-backend-api/api/domain/entities"
	"github.com/thalesmacedo1/freterapido-backend-api/api/interfaces/api"
)

// Mock for the GetShippingQuotationUseCase
type MockGetShippingQuotationUseCase struct {
	mock.Mock
}

func (m *MockGetShippingQuotationUseCase) Execute(ctx gin.Context, request domain.QuoteRequest) (*domain.QuoteResponse, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.QuoteResponse), args.Error(1)
}

func TestQuoteController_GetQuote_Success(t *testing.T) {
	// Setup Gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	// Setup mock use case
	mockUseCase := new(MockGetShippingQuotationUseCase)

	// Create controller
	controller := api.NewQuoteController(mockUseCase)

	// Setup request
	requestBody := domain.QuoteRequest{}
	requestBody.Recipient.Address.Zipcode = "01311000"

	volume := domain.Volume{
		Category:      7,
		Amount:        1,
		UnitaryWeight: 5.0,
		Price:         349.0,
		SKU:           "abc-teste-123",
		Height:        0.2,
		Width:         0.2,
		Length:        0.2,
	}

	requestBody.Volumes = append(requestBody.Volumes, volume)

	// Setup expected response
	expectedResponse := &domain.QuoteResponse{}
	carrier := domain.Carrier{
		Name:     "EXPRESSO FR",
		Service:  "Rodoviário",
		Deadline: "3",
		Price:    17.0,
	}
	expectedResponse.Carriers = append(expectedResponse.Carriers, carrier)

	// Setup mock expectations (using any request since we'd need to match the exact same struct)
	mockUseCase.On("Execute", mock.Anything, mock.AnythingOfType("domain.QuoteRequest")).Return(expectedResponse, nil)

	// Convert request to JSON
	jsonBody, _ := json.Marshal(requestBody)

	// Create request
	req, _ := http.NewRequest("POST", "/quote", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Setup router
	r.POST("/quote", controller.GetQuote)

	// Serve request
	r.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse response body
	var responseBody domain.QuoteResponse
	json.Unmarshal(w.Body.Bytes(), &responseBody)

	assert.Len(t, responseBody.Carriers, 1)
	assert.Equal(t, "EXPRESSO FR", responseBody.Carriers[0].Name)
	assert.Equal(t, "Rodoviário", responseBody.Carriers[0].Service)
	assert.Equal(t, "3", responseBody.Carriers[0].Deadline)
	assert.Equal(t, 17.0, responseBody.Carriers[0].Price)

	// Verify all expectations were met
	mockUseCase.AssertExpectations(t)
}

func TestQuoteController_GetQuote_ValidationError(t *testing.T) {
	// Setup Gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	// Setup mock use case
	mockUseCase := new(MockGetShippingQuotationUseCase)

	// Create controller
	controller := api.NewQuoteController(mockUseCase)

	// Setup invalid request (missing zipcode)
	requestBody := domain.QuoteRequest{}

	volume := domain.Volume{
		Category:      7,
		Amount:        1,
		UnitaryWeight: 5.0,
		Price:         349.0,
		SKU:           "abc-teste-123",
		Height:        0.2,
		Width:         0.2,
		Length:        0.2,
	}

	requestBody.Volumes = append(requestBody.Volumes, volume)

	// Convert request to JSON
	jsonBody, _ := json.Marshal(requestBody)

	// Create request
	req, _ := http.NewRequest("POST", "/quote", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Setup router
	r.POST("/quote", controller.GetQuote)

	// Serve request
	r.ServeHTTP(w, req)

	// Assert response - should be bad request due to missing zipcode
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// The mock should not have been called
	mockUseCase.AssertNotCalled(t, "Execute")
}

func TestQuoteController_GetQuote_UseCaseError(t *testing.T) {
	// Setup Gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	// Setup mock use case
	mockUseCase := new(MockGetShippingQuotationUseCase)

	// Create controller
	controller := api.NewQuoteController(mockUseCase)

	// Setup request
	requestBody := domain.QuoteRequest{}
	requestBody.Recipient.Address.Zipcode = "01311000"

	volume := domain.Volume{
		Category:      7,
		Amount:        1,
		UnitaryWeight: 5.0,
		Price:         349.0,
		SKU:           "abc-teste-123",
		Height:        0.2,
		Width:         0.2,
		Length:        0.2,
	}

	requestBody.Volumes = append(requestBody.Volumes, volume)

	// Setup mock to return error
	expectedError := errors.New("use case error")
	mockUseCase.On("Execute", mock.Anything, mock.AnythingOfType("domain.QuoteRequest")).Return(nil, expectedError)

	// Convert request to JSON
	jsonBody, _ := json.Marshal(requestBody)

	// Create request
	req, _ := http.NewRequest("POST", "/quote", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Setup router
	r.POST("/quote", controller.GetQuote)

	// Serve request
	r.ServeHTTP(w, req)

	// Assert response - should be internal server error
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// Verify expectations were met
	mockUseCase.AssertExpectations(t)
}
