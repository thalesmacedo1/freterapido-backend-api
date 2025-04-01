package usecases

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	domain "github.com/thalesmacedo1/freterapido-backend-api/api/domain/entities"
)

type GetShippingQuotationUseCase struct {
	quoteRepository domain.QuoteRepository
}

func NewGetShippingQuotationUseCase(quoteRepository domain.QuoteRepository) *GetShippingQuotationUseCase {
	return &GetShippingQuotationUseCase{
		quoteRepository: quoteRepository,
	}
}

func (uc *GetShippingQuotationUseCase) Execute(ctx context.Context, request domain.QuoteRequest) (*domain.QuoteResponse, error) {
	// Prepare request to Frete Rápido API
	freteRapidoRequest := prepareFRRequest(request)

	// Call Frete Rápido API
	frResponse, err := callFreteRapidoAPI(freteRapidoRequest)
	if err != nil {
		return nil, fmt.Errorf("error calling Frete Rápido API: %w", err)
	}

	// Transform response
	quoteResponse := transformResponse(frResponse)

	// Save quote to database
	err = uc.quoteRepository.SaveQuote(ctx, quoteResponse)
	if err != nil {
		return nil, fmt.Errorf("error saving quote: %w", err)
	}

	return quoteResponse, nil
}

func prepareFRRequest(request domain.QuoteRequest) domain.FreteRapidoRequest {
	frRequest := domain.FreteRapidoRequest{}

	// Set shipper information using environment variables
	frRequest.Shipper.RegisteredNumber = os.Getenv("SHIPPER_CNPJ")
	frRequest.Shipper.Token = os.Getenv("FRETE_RAPIDO_TOKEN")
	frRequest.Shipper.PlatformCode = os.Getenv("PLATFORM_CODE")

	// Set recipient information
	frRequest.Recipient.Address.Zipcode = request.Recipient.Address.Zipcode

	// Create dispatcher
	dispatcher := struct {
		RegisteredNumber string                     `json:"registered_number"`
		Zipcode          string                     `json:"zipcode"`
		Volumes          []domain.FreteRapidoVolume `json:"volumes"`
	}{
		RegisteredNumber: os.Getenv("SHIPPER_CNPJ"),
		Zipcode:          os.Getenv("DISPATCHER_ZIPCODE"),
		Volumes:          []domain.FreteRapidoVolume{},
	}

	// Transform volumes
	for _, vol := range request.Volumes {
		frVolume := domain.FreteRapidoVolume{
			Amount:        vol.Amount,
			Category:      vol.Category,
			SKU:           vol.SKU,
			Height:        vol.Height,
			Width:         vol.Width,
			Length:        vol.Length,
			UnitaryWeight: vol.UnitaryWeight,
			UnitaryPrice:  vol.Price,
		}
		dispatcher.Volumes = append(dispatcher.Volumes, frVolume)
	}

	frRequest.Dispatchers = append(frRequest.Dispatchers, dispatcher)

	return frRequest
}

func callFreteRapidoAPI(request domain.FreteRapidoRequest) (*domain.FreteRapidoResponse, error) {
	// Convert request to JSON
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	// Create HTTP request using environment variable for API URL
	req, err := http.NewRequest("POST", os.Getenv("FRETE_RAPIDO_API_URL"), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response status: %d", resp.StatusCode)
	}

	// Parse response
	var frResponse domain.FreteRapidoResponse
	if err := json.NewDecoder(resp.Body).Decode(&frResponse); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &frResponse, nil
}

func transformResponse(frResponse *domain.FreteRapidoResponse) *domain.QuoteResponse {
	response := &domain.QuoteResponse{
		Carriers: []domain.Carrier{},
	}

	// Extract carriers from response
	for _, dispatcher := range frResponse.Dispatchers {
		for _, offer := range dispatcher.Offers {
			carrier := domain.Carrier{
				Name:     offer.Carrier.Name,
				Service:  offer.Service,
				Deadline: strconv.Itoa(offer.DeliveryTime.Days),
				Price:    offer.FinalPrice,
			}
			response.Carriers = append(response.Carriers, carrier)
		}
	}

	return response
}
