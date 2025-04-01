package usecases

import (
	"context"

	domain "github.com/thalesmacedo1/freterapido-backend-api/api/domain/entities"
)

type GetLastQuotesUseCase struct {
	quoteRepository domain.QuoteRepository
}

func NewGetLastQuotesUseCase(quoteRepository domain.QuoteRepository) *GetLastQuotesUseCase {
	return &GetLastQuotesUseCase{
		quoteRepository: quoteRepository,
	}
}

func (uc *GetLastQuotesUseCase) Execute(ctx context.Context, limit int) ([]domain.QuoteResponse, error) {
	// Get last quotes from repository
	return uc.quoteRepository.GetLastQuotes(ctx, limit)
}
