package services

import (
	"context"

	"github.com/afikrim/go-hexa-template/internal/core/domains"
)

type CountryService interface {
	FindAll(ctx context.Context) ([]domains.Country, error)
}
