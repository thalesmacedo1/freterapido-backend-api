package domain

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type QuoteRequest struct {
	Recipient struct {
		Address struct {
			Zipcode string `json:"zipcode"`
		} `json:"address"`
	} `json:"recipient"`
	Volumes []Volume `json:"volumes"`
}

type Volume struct {
	Category      int     `json:"category"`
	Amount        int     `json:"amount"`
	UnitaryWeight float64 `json:"unitary_weight"`
	Price         float64 `json:"price"`
	SKU           string  `json:"sku"`
	Height        float64 `json:"height"`
	Width         float64 `json:"width"`
	Length        float64 `json:"length"`
}

type QuoteResponse struct {
	gorm.Model
	Carriers []Carrier `json:"carrier" gorm:"type:jsonb"`
}

type Carrier struct {
	Name     string  `json:"name"`
	Service  string  `json:"service"`
	Deadline string  `json:"deadline"`
	Price    float64 `json:"price"`
}

type QuoteRepository interface {
	GetTotalCasesAndDeaths(ctx context.Context, countryCode string, date time.Time) (*valueobjects.Quote, error)
}
