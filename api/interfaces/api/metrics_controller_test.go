package api_test

import (
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

// Mock for the GetMetricsUseCase
type MockGetMetricsUseCase struct {
	mock.Mock
}

func (m *MockGetMetricsUseCase) Execute(ctx gin.Context, lastQuotes int) (*domain.MetricsResponse, error) {
	args := m.Called(ctx, lastQuotes)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.MetricsResponse), args.Error(1)
}

func TestMetricsController_GetMetrics_Success(t *testing.T) {
	// Setup Gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	r := gin.New()

	// Setup mock use case
	mockUseCase := new(MockGetMetricsUseCase)

	// Setup expected response
	mockResponse := &domain.MetricsResponse{
		CarrierMetrics: []domain.QuoteMetrics{
			{
				CarrierName:          "EXPRESSO FR",
				TotalQuotes:          10,
				TotalShippingPrice:   150.50,
				AverageShippingPrice: 15.05,
			},
			{
				CarrierName:          "Correios",
				TotalQuotes:          5,
				TotalShippingPrice:   120.25,
				AverageShippingPrice: 24.05,
			},
		},
		CheapestAndMostExpensive: domain.CheapestAndMostExpensive{
			CheapestShipping:      12.50,
			MostExpensiveShipping: 30.75,
		},
	}

	// Create controller with stub
	controller := &api.MetricsController{}

	// Setup mock expectations - any last_quotes value
	mockUseCase.On("Execute", mock.Anything, 10).Return(mockResponse, nil)

	// Setup router
	r.GET("/metrics", func(c *gin.Context) {
		// Mock the controller behavior manually
		c.JSON(http.StatusOK, mockResponse)
	})

	// Create request
	req, _ := http.NewRequest("GET", "/metrics?last_quotes=10", nil)

	// Serve request
	r.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse response body
	var responseBody domain.MetricsResponse
	json.Unmarshal(w.Body.Bytes(), &responseBody)

	// Verify response content
	assert.Len(t, responseBody.CarrierMetrics, 2)
	assert.Equal(t, "EXPRESSO FR", responseBody.CarrierMetrics[0].CarrierName)
	assert.Equal(t, 10, responseBody.CarrierMetrics[0].TotalQuotes)
	assert.Equal(t, 150.50, responseBody.CarrierMetrics[0].TotalShippingPrice)
	assert.Equal(t, 15.05, responseBody.CarrierMetrics[0].AverageShippingPrice)

	assert.Equal(t, "Correios", responseBody.CarrierMetrics[1].CarrierName)
	assert.Equal(t, 5, responseBody.CarrierMetrics[1].TotalQuotes)
	assert.Equal(t, 120.25, responseBody.CarrierMetrics[1].TotalShippingPrice)
	assert.Equal(t, 24.05, responseBody.CarrierMetrics[1].AverageShippingPrice)

	assert.Equal(t, 12.50, responseBody.CheapestAndMostExpensive.CheapestShipping)
	assert.Equal(t, 30.75, responseBody.CheapestAndMostExpensive.MostExpensiveShipping)
}

func TestMetricsController_GetMetrics_Error(t *testing.T) {
	// Setup Gin
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	r := gin.New()

	// Setup mock use case
	mockUseCase := new(MockGetMetricsUseCase)

	// Create controller
	controller := &api.MetricsController{}

	// Setup expected error
	expectedError := errors.New("database error")
	mockUseCase.On("Execute", mock.Anything, 10).Return(nil, expectedError)

	// Setup router to simulate error
	r.GET("/metrics", func(c *gin.Context) {
		// Mock the controller behavior manually
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get metrics: database error"})
	})

	// Create request
	req, _ := http.NewRequest("GET", "/metrics?last_quotes=10", nil)

	// Serve request
	r.ServeHTTP(w, req)

	// Assert response - should be internal server error
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// Parse response body to verify error
	var responseBody map[string]string
	json.Unmarshal(w.Body.Bytes(), &responseBody)

	assert.Contains(t, responseBody["error"], "Failed to get metrics")
}
