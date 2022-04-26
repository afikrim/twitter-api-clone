package country_service

import (
	"context"

	"github.com/afikrim/go-hexa-template/internal/core/domains"
	"github.com/afikrim/go-hexa-template/internal/core/ports/repositories"
)

type service struct {
	repo repositories.CountryRepository
}

func NewCountryService(repo repositories.CountryRepository) *service {
	return &service{
		repo: repo,
	}
}

func (s *service) FindAll(ctx context.Context) ([]domains.Country, error) {
	return s.repo.FindAll(ctx)
}
