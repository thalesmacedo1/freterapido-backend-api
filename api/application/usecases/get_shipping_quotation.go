package usecases

import (
	"context"
	"time"

	"github.com/thalesmacedo1/covid-api/domain/entities"
	"github.com/thalesmacedo1/covid-api/domain/repositories"
)

type GetCountryWithMostCasesUseCase interface {
	Execute(ctx context.Context, input GetCountryWithMostCasesInput) (*GetCountryWithMostCasesOutput, error)
}

type GetCountryWithMostCasesInput struct {
	Date time.Time
}

type GetCountryWithMostCasesOutput struct {
	Country         entities.Country
	CumulativeCases int
}

type getCountryWithMostCasesUseCase struct {
	covidStatsRepo repositories.CovidStatsRepository
	countryRepo    repositories.CountryRepository
}

func NewGetCountryWithMostCasesUseCase(covidRepo repositories.CovidStatsRepository, countryRepo repositories.CountryRepository) GetCountryWithMostCasesUseCase {
	return &getCountryWithMostCasesUseCase{
		covidStatsRepo: covidRepo,
		countryRepo:    countryRepo,
	}
}

func (uc *getCountryWithMostCasesUseCase) Execute(ctx context.Context, input GetCountryWithMostCasesInput) (*GetCountryWithMostCasesOutput, error) {

	countryWithMostCases, cumulativeCases, err := uc.covidStatsRepo.GetCountryWithMostCases(ctx, input.Date)
	if err != nil {
		return nil, err
	}

	country, err := uc.countryRepo.GetCountryByCode(ctx, countryWithMostCases)
	if err != nil {
		return nil, err
	}

	return &GetCountryWithMostCasesOutput{
		Country:         *country,
		CumulativeCases: cumulativeCases,
	}, nil
}
