package country_repository

import "github.com/afikrim/go-hexa-template/internal/core/domains"

type Country struct {
	ID   uint64 `gorm:"column:id;not null;primaryKey;autoIncrement"`
	Code string `gorm:"column:code;type:varchar(255);not null;unique"`
	Name string `gorm:"column:name;type:varchar(255);not null"`
}

func (c *Country) TableName() string {
	return "countries"
}

func (c *Country) ToDomain() *domains.Country {
	return &domains.Country{
		ID:   c.ID,
		Code: c.Code,
		Name: c.Name,
	}
}
