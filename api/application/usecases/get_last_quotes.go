package usecases

import (
	"context"
	"time"

	"github.com/thalesmacedo1/covid-api/domain/entities"
	"github.com/thalesmacedo1/covid-api/domain/repositories"
)

type GetLastQuotesUseCase interface {
	Execute(ctx context.Context, input GetLastQuotesInput) (*GetLastQuotesOutput, error)
}

type GetLastQuotesInput struct {
	Date time.Time
}

type GetLastQuotesOutput struct {
	Country         entities.Country
	CumulativeCases int
}

type getLastQuotesUseCase struct {
	countryRepo repositories.CountryRepository
}

func NewGetLastQuotesUseCase(covidRepo repositories.CovidStatsRepository, countryRepo repositories.CountryRepository) GetLastQuotesUseCase {
	return &getLastQuotesUseCase{
		countryRepo: countryRepo,
	}
}

func (uc *getLastQuotesUseCase) Execute(ctx context.Context, input GetLastQuotesInput) (*GetLastQuotesOutput, error) {

	countryWithMostCases, cumulativeCases, err := uc.covidStatsRepo.GetLastQuotes(ctx, input.Date)
	if err != nil {
		return nil, err
	}

	country, err := uc.countryRepo.GetCountryByCode(ctx, countryWithMostCases)
	if err != nil {
		return nil, err
	}

	return &GetLastQuotesOutput{
		Country: *country,
	}, nil
}
