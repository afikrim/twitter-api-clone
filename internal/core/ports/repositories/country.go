package repositories

import (
	"context"

	"github.com/afikrim/go-hexa-template/internal/core/domains"
)

type CountryRepository interface {
	FindAll(ctx context.Context) ([]domains.Country, error)
}
