package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	domain "github.com/thalesmacedo1/freterapido-backend-api/api/domain/entities"
)

func TestQuoteEndpoint_Integration(t *testing.T) {
	// Skip if test environment is not set up
	if testRouter == nil {
		t.Skip("Test environment not set up")
	}

	// Create test request body
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

	// Convert to JSON
	jsonBody, err := json.Marshal(requestBody)
	assert.NoError(t, err)

	// Create HTTP request
	req, err := http.NewRequest("POST", "/quote", bytes.NewBuffer(jsonBody))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	testRouter.ServeHTTP(w, req)

	// Check response status
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse response body
	var response domain.QuoteResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Check response content
	assert.NotEmpty(t, response.Carriers)

	// Validate that quotes were saved to the database
	quotes, err := testQuoteRepository.GetLastQuotes(req.Context(), 10)
	assert.NoError(t, err)
	assert.NotEmpty(t, quotes)
}
