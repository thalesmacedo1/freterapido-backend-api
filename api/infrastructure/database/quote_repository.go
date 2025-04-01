package database

import (
	"context"

	domain "github.com/thalesmacedo1/freterapido-backend-api/api/domain/entities"
	"gorm.io/gorm"
)

type QuoteRepositoryImpl struct {
	db *gorm.DB
}

func NewQuoteRepository(db *gorm.DB) domain.QuoteRepository {
	return &QuoteRepositoryImpl{
		db: db,
	}
}

func (r *QuoteRepositoryImpl) SaveQuote(ctx context.Context, quote *domain.QuoteResponse) error {
	result := r.db.Create(quote)
	return result.Error
}

func (r *QuoteRepositoryImpl) GetLastQuotes(ctx context.Context, limit int) ([]domain.QuoteResponse, error) {
	var quotes []domain.QuoteResponse
	
	query := r.db.Order("created_at DESC")
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	
	result := query.Find(&quotes)
	
	if result.Error != nil {
		return nil, result.Error
	}
	
	return quotes, nil
} 