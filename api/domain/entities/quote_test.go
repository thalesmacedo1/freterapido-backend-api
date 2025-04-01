package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	domain "github.com/thalesmacedo1/freterapido-backend-api/api/domain/entities"
)

func TestQuoteRequest(t *testing.T) {
	// Create a valid quote request
	request := domain.QuoteRequest{}
	request.Recipient.Address.Zipcode = "01311000"

	volume1 := domain.Volume{
		Category:      7,
		Amount:        1,
		UnitaryWeight: 5.0,
		Price:         349.0,
		SKU:           "abc-teste-123",
		Height:        0.2,
		Width:         0.2,
		Length:        0.2,
	}

	volume2 := domain.Volume{
		Category:      7,
		Amount:        2,
		UnitaryWeight: 4.0,
		Price:         556.0,
		SKU:           "abc-teste-527",
		Height:        0.4,
		Width:         0.6,
		Length:        0.15,
	}

	request.Volumes = append(request.Volumes, volume1, volume2)

	// Assert the values are set correctly
	assert.Equal(t, "01311000", request.Recipient.Address.Zipcode)
	assert.Len(t, request.Volumes, 2)
	assert.Equal(t, 7, request.Volumes[0].Category)
	assert.Equal(t, 1, request.Volumes[0].Amount)
	assert.Equal(t, 5.0, request.Volumes[0].UnitaryWeight)
	assert.Equal(t, 349.0, request.Volumes[0].Price)
	assert.Equal(t, "abc-teste-123", request.Volumes[0].SKU)
	assert.Equal(t, 0.2, request.Volumes[0].Height)
	assert.Equal(t, 0.2, request.Volumes[0].Width)
	assert.Equal(t, 0.2, request.Volumes[0].Length)
}

func TestQuoteResponse(t *testing.T) {
	// Create a quote response
	response := domain.QuoteResponse{}

	carrier1 := domain.Carrier{
		Name:     "EXPRESSO FR",
		Service:  "Rodoviário",
		Deadline: "3",
		Price:    17.0,
	}

	carrier2 := domain.Carrier{
		Name:     "Correios",
		Service:  "SEDEX",
		Deadline: "1",
		Price:    20.99,
	}

	response.Carriers = append(response.Carriers, carrier1, carrier2)

	// Assert the values are set correctly
	assert.Len(t, response.Carriers, 2)
	assert.Equal(t, "EXPRESSO FR", response.Carriers[0].Name)
	assert.Equal(t, "Rodoviário", response.Carriers[0].Service)
	assert.Equal(t, "3", response.Carriers[0].Deadline)
	assert.Equal(t, 17.0, response.Carriers[0].Price)

	assert.Equal(t, "Correios", response.Carriers[1].Name)
	assert.Equal(t, "SEDEX", response.Carriers[1].Service)
	assert.Equal(t, "1", response.Carriers[1].Deadline)
	assert.Equal(t, 20.99, response.Carriers[1].Price)
}
