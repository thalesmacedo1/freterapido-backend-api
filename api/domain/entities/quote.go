package domain

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Input structures (from clients to our API)
// QuoteRequest representa a solicitação de cotação de frete
// @Description Solicitação para obter cotações de frete de diferentes transportadoras
type QuoteRequest struct {
	// Informações do destinatário
	Recipient struct {
		// Endereço do destinatário
		Address struct {
			// CEP do destinatário (obrigatório)
			// @example "01311000"
			Zipcode string `json:"zipcode"`
		} `json:"address"`
	} `json:"recipient"`
	// Lista de volumes para transporte
	// @Description Lista de volumes para cálculo de frete
	// @Required
	Volumes []Volume `json:"volumes"`
}

// Volume representa um volume para transporte
// @Description Detalhes de um volume para cotação de frete
type Volume struct {
	// Categoria do produto
	// @example 7
	Category int `json:"category"`
	// Quantidade de itens
	// @example 1
	Amount int `json:"amount"`
	// Peso unitário em kg
	// @example 5.0
	UnitaryWeight float64 `json:"unitary_weight"`
	// Preço unitário do produto
	// @example 349.90
	Price float64 `json:"price"`
	// Código SKU do produto
	// @example "abc-teste-123"
	SKU string `json:"sku"`
	// Altura do volume em metros
	// @example 0.2
	Height float64 `json:"height"`
	// Largura do volume em metros
	// @example 0.2
	Width float64 `json:"width"`
	// Comprimento do volume em metros
	// @example 0.2
	Length float64 `json:"length"`
}

// API response structure (from our API to clients)
// QuoteResponse representa a resposta com cotações de frete
// @Description Resposta com as cotações de frete disponíveis
type QuoteResponse struct {
	gorm.Model
	// Lista de transportadoras com suas cotações
	// @Description Lista de transportadoras e seus valores
	Carriers CarriersJSON `json:"carrier" gorm:"type:jsonb"`
}

// CarriersJSON é um tipo personalizado para serializar como JSONB no PostgreSQL
type CarriersJSON []Carrier

// Implementação da interface driver.Valuer
func (c CarriersJSON) Value() (driver.Value, error) {
	if len(c) == 0 {
		return nil, nil
	}
	return json.Marshal(c)
}

// Implementação da interface sql.Scanner
func (c *CarriersJSON) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	var data []byte
	switch v := value.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return fmt.Errorf("cannot scan %T into CarriersJSON", value)
	}
	return json.Unmarshal(data, c)
}

// Carrier representa as informações de uma transportadora
// @Description Informações sobre a cotação de uma transportadora específica
type Carrier struct {
	// Nome da transportadora
	// @example "EXPRESSO FR"
	Name string `json:"name"`
	// Serviço oferecido
	// @example "Rodoviário"
	Service string `json:"service"`
	// Prazo de entrega em dias
	// @example "3"
	Deadline string `json:"deadline"`
	// Valor do frete
	// @example 17.00
	Price float64 `json:"price"`
}

// Frete Rápido API structure
type FreteRapidoRequest struct {
	Shipper struct {
		RegisteredNumber string `json:"registered_number"`
		Token            string `json:"token"`
		PlatformCode     string `json:"platform_code"`
	} `json:"shipper"`
	Recipient struct {
		Type    int    `json:"type"`
		Country string `json:"country"`
		Zipcode int    `json:"zipcode"`
	} `json:"recipient"`
	Dispatchers []struct {
		RegisteredNumber string              `json:"registered_number"`
		Zipcode          int                 `json:"zipcode"`
		Volumes          []FreteRapidoVolume `json:"volumes"`
	} `json:"dispatchers"`
	SimulationType []int `json:"simulation_type"`
	Returns        struct {
		Composition  bool `json:"composition"`
		Volumes      bool `json:"volumes"`
		AppliedRules bool `json:"applied_rules"`
	} `json:"returns"`
}

type FreteRapidoVolume struct {
	Amount        int     `json:"amount"`
	Category      string  `json:"category"`
	Sku           string  `json:"sku"`
	Height        float64 `json:"height"`
	Width         float64 `json:"width"`
	Length        float64 `json:"length"`
	UnitaryPrice  float64 `json:"unitary_price"`
	UnitaryWeight float64 `json:"unitary_weight"`
}

type FreteRapidoResponse struct {
	Dispatchers []struct {
		ID                         string `json:"id"`
		RequestID                  string `json:"request_id"`
		RegisteredNumberShipper    string `json:"registered_number_shipper"`
		RegisteredNumberDispatcher string `json:"registered_number_dispatcher"`
		ZipcodeOrigin              int    `json:"zipcode_origin"`
		Offers                     []struct {
			Offer          int    `json:"offer"`
			TableReference string `json:"table_reference"`
			SimulationType int    `json:"simulation_type"`
			Carrier        struct {
				Name             string `json:"name"`
				RegisteredNumber string `json:"registered_number"`
				StateInscription string `json:"state_inscription"`
				Logo             string `json:"logo"`
				Reference        int    `json:"reference"`
				CompanyName      string `json:"company_name"`
			} `json:"carrier"`
			Service      string `json:"service"`
			DeliveryTime struct {
				Days          int    `json:"days"`
				EstimatedDate string `json:"estimated_date"`
			} `json:"delivery_time"`
			Expiration time.Time `json:"expiration"`
			CostPrice  float64   `json:"cost_price"`
			FinalPrice float64   `json:"final_price"`
			Weights    struct {
				Real  float64 `json:"real"`
				Cubed float64 `json:"cubed"`
				Used  float64 `json:"used"`
			} `json:"weights"`
			OriginalDeliveryTime struct {
				Days          int    `json:"days"`
				EstimatedDate string `json:"estimated_date"`
			} `json:"original_delivery_time"`
			HomeDelivery                bool `json:"home_delivery"`
			CarrierOriginalDeliveryTime struct {
				Days          int    `json:"days"`
				EstimatedDate string `json:"estimated_date"`
			} `json:"carrier_original_delivery_time"`
			Modal string `json:"modal"`
		} `json:"offers"`
	} `json:"dispatchers"`
}

// Repository interface
type QuoteRepository interface {
	SaveQuote(ctx context.Context, quote *QuoteResponse) error
	GetLastQuotes(ctx context.Context, limit int) ([]QuoteResponse, error)
}
