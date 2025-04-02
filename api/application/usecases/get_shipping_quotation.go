package usecases

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
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

// Função auxiliar para obter variáveis de ambiente com valor padrão
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Printf("WARNING: Environment variable %s not found, using default value: %s", key, defaultValue)
		return defaultValue
	}
	return value
}

func (uc *GetShippingQuotationUseCase) Execute(ctx context.Context, request domain.QuoteRequest) (*domain.QuoteResponse, error) {
	freteRapidoRequest := prepareFRRequest(request)

	frResponse, err := callFreteRapidoAPI(freteRapidoRequest)
	if err != nil {
		return nil, fmt.Errorf("error calling Frete Rápido API: %w", err)
	}

	quoteResponse := transformResponse(frResponse)

	err = uc.quoteRepository.SaveQuote(ctx, quoteResponse)
	if err != nil {
		return nil, fmt.Errorf("error saving quote: %w", err)
	}

	return quoteResponse, nil
}

func prepareFRRequest(request domain.QuoteRequest) domain.FreteRapidoRequest {
	frRequest := domain.FreteRapidoRequest{}

	frRequest.Shipper.RegisteredNumber = getEnv("CNPJ", "25438296000158")
	frRequest.Shipper.Token = getEnv("FRETE_RAPIDO_TOKEN", "1d52a9b6b78cf07b08586152459a5c90")
	frRequest.Shipper.PlatformCode = getEnv("PLATFORM_CODE", "5AKVkHqCn")

	frRequest.Recipient.Type = 0
	frRequest.Recipient.Country = "BRA"

	zipcodeInt, _ := strconv.Atoi(request.Recipient.Address.Zipcode)
	frRequest.Recipient.Zipcode = zipcodeInt

	zipcodeEnvInt, _ := strconv.Atoi(getEnv("ZIPCODE", "29161376"))
	dispatcher := struct {
		RegisteredNumber string                     `json:"registered_number"`
		Zipcode          int                        `json:"zipcode"`
		Volumes          []domain.FreteRapidoVolume `json:"volumes"`
	}{
		RegisteredNumber: getEnv("CNPJ", "25438296000158"),
		Zipcode:          zipcodeEnvInt,
		Volumes:          []domain.FreteRapidoVolume{},
	}

	for _, vol := range request.Volumes {
		frVolume := domain.FreteRapidoVolume{
			Amount:        vol.Amount,
			Category:      strconv.Itoa(vol.Category),
			Sku:           vol.SKU,
			Height:        vol.Height,
			Width:         vol.Width,
			Length:        vol.Length,
			UnitaryWeight: vol.UnitaryWeight,
			UnitaryPrice:  vol.Price,
		}
		dispatcher.Volumes = append(dispatcher.Volumes, frVolume)
	}

	frRequest.Dispatchers = append(frRequest.Dispatchers, dispatcher)
	frRequest.SimulationType = []int{0}
	frRequest.Returns.Composition = false
	frRequest.Returns.Volumes = false
	frRequest.Returns.AppliedRules = false

	jsonData, _ := json.MarshalIndent(frRequest, "", "  ")
	log.Printf("FreteRapido Request: %s", string(jsonData))

	return frRequest
}

func callFreteRapidoAPI(request domain.FreteRapidoRequest) (*domain.FreteRapidoResponse, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	apiURL := getEnv("FRETE_RAPIDO_API_URL", "https://sp.freterapido.com/api/v3/quote/simulate")
	log.Printf("Calling FreteRapido API at: %s", apiURL)

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody := make([]byte, 1024)
		n, _ := resp.Body.Read(respBody)
		return nil, fmt.Errorf("received non-200 response status: %d, body: %s", resp.StatusCode, string(respBody[:n]))
	}

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
