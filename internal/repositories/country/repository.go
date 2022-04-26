package country_repository

import (
	"context"

	"github.com/afikrim/go-hexa-template/internal/core/domains"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewCountryRepository(db *gorm.DB) *repository {
	return &repository{
		db: db,
	}
}

func (r *repository) FindAll(ctx context.Context) ([]domains.Country, error) {
	var countryModel []Country
	if err := r.db.Find(&countryModel).WithContext(ctx).Error; err != nil {
		return nil, err
	}

	var countries []domains.Country
	for _, country := range countryModel {
		countries = append(countries, *country.ToDomain())
	}

	return countries, nil
}
