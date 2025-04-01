package domain

import (
	"context"

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
	Carriers []Carrier `json:"carrier" gorm:"type:jsonb"`
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

// Frete Rápido API structures (for external API call)
type FreteRapidoRequest struct {
	Shipper struct {
		RegisteredNumber string `json:"registered_number"`
		Token            string `json:"token"`
		PlatformCode     string `json:"platform_code"`
	} `json:"shipper"`
	Dispatchers []struct {
		RegisteredNumber string              `json:"registered_number"`
		Zipcode          string              `json:"zipcode"`
		Volumes          []FreteRapidoVolume `json:"volumes"`
	} `json:"dispatchers"`
	Recipient struct {
		Address struct {
			Zipcode string `json:"zipcode"`
		} `json:"address"`
	} `json:"recipient"`
}

type FreteRapidoVolume struct {
	Amount        int     `json:"amount"`
	Category      int     `json:"category"`
	SKU           string  `json:"sku"`
	Height        float64 `json:"height"`
	Width         float64 `json:"width"`
	Length        float64 `json:"length"`
	UnitaryWeight float64 `json:"unitary_weight"`
	UnitaryPrice  float64 `json:"unitary_price"`
}

type FreteRapidoResponse struct {
	Dispatchers []struct {
		Offers []struct {
			Carrier struct {
				Name string `json:"name"`
			} `json:"carrier"`
			Service      string  `json:"service"`
			FinalPrice   float64 `json:"final_price"`
			DeliveryTime struct {
				Days int `json:"days"`
			} `json:"delivery_time"`
		} `json:"offers"`
	} `json:"dispatchers"`
}

// Repository interface
type QuoteRepository interface {
	SaveQuote(ctx context.Context, quote *QuoteResponse) error
	GetLastQuotes(ctx context.Context, limit int) ([]QuoteResponse, error)
}
