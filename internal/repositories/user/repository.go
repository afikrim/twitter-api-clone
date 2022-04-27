package user_repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/afikrim/go-hexa-template/internal/core/domains"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(ctx context.Context, dto *domains.RegisterDto) (*domains.User, error) {
	userModel := User{}.FromRegisterDto(dto)
	if err := r.db.WithContext(ctx).Create(&userModel).Error; err != nil {
		return nil, err
	}

	return userModel.ToDomain(), nil
}

func (r *repository) FindAll(ctx context.Context, query *domains.QueryParamDto) ([]domains.User, error) {
	dbQuery := r.db.WithContext(ctx)
	if query.Search != "" {
		search := fmt.Sprintf("%%%s%%", query.Search)
		dbQuery = dbQuery.Where("username LIKE ? OR fullname LIKE ?", search, search)
	}

	if query.SortBy != "" && query.OrderBy != "" {
		orderby := fmt.Sprintf("%s %s", query.SortBy, query.OrderBy)
		dbQuery = dbQuery.Order(orderby)
	}

	if query.Limit != nil {
		dbQuery = dbQuery.Limit(*query.Limit)
	}

	if query.Offset != nil {
		dbQuery = dbQuery.Offset(*query.Offset)
	}

	var userModels []User
	if err := dbQuery.Find(&userModels).Error; err != nil {
		return nil, err
	}

	var users []domains.User
	for _, userModel := range userModels {
		users = append(users, *userModel.ToDomain())
	}

	return users, nil
}

func (r *repository) FindByID(ctx context.Context, id int64) (*domains.User, error) {
	var userModel User
	if err := r.db.WithContext(ctx).First(&userModel, id).Error; err != nil {
		return nil, err
	}

	return userModel.ToDomainWithTimestamps(), nil
}

func (r *repository) FindByUsername(ctx context.Context, username string) (*domains.User, error) {
	var userModel User
	if err := r.db.WithContext(ctx).Joins("Country").Where("username = ?", username).First(&userModel).Error; err != nil {
		return nil, err
	}

	return userModel.ToDomainWithCountryAndTimestamps(), nil
}

func (r *repository) FindByEmail(ctx context.Context, email string) (*domains.User, error) {
	var userModel User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&userModel).Error; err != nil {
		return nil, err
	}

	return userModel.ToDomainWithTimestamps(), nil
}

func (r *repository) FindByPhone(ctx context.Context, phone string) (*domains.User, error) {
	var userModel User
	if err := r.db.WithContext(ctx).Where("phone = ?", phone).First(&userModel).Error; err != nil {
		return nil, err
	}

	return userModel.ToDomainWithTimestamps(), nil
}

func (r *repository) Update(ctx context.Context, id int64, dto *domains.UpdateUserDto) (*domains.User, error) {
	var userModel User
	if err := r.db.WithContext(ctx).Joins("Country").First(&userModel, id).Error; err != nil {
		return nil, err
	}

	if dto.Fullname != "" {
		userModel.Fullname = dto.Fullname
	}

	if dto.Gender != nil {
		userModel.Gender = *dto.Gender
	}

	if dto.BirthDate != "" {
		parsedBirthDate, err := time.Parse("2006-01-02", dto.BirthDate)
		if err != nil {
			return nil, errors.New("Birthdate format is wrong")
		}

		userModel.BirthDate = parsedBirthDate
	}

	if dto.CountryID != 0 {
		userModel.CountryID = dto.CountryID
	}

	if err := r.db.WithContext(ctx).Save(&userModel).Error; err != nil {
		return nil, err
	}

	return userModel.ToDomainWithCountryAndTimestamps(), nil
}

func (r *repository) UpdateCredential(ctx context.Context, id int64, dto *domains.UpdateUserCredentialDto) (*domains.User, error) {
	var userModel User
	if err := r.db.WithContext(ctx).Joins("Country").First(&userModel, id).Error; err != nil {
		return nil, err
	}

	if dto.Username != "" {
		userModel.Username = dto.Username
	}

	if dto.Phone != "" {
		userModel.Phone = dto.Phone
	}

	if dto.Email != "" {
		userModel.Email = dto.Email
	}

	if err := r.db.WithContext(ctx).Save(&userModel).Error; err != nil {
		return nil, err
	}

	return userModel.ToDomainWithCountryAndTimestamps(), nil
}

func (r *repository) UpdatePassword(ctx context.Context, id int64, dto *domains.UpdateUserPasswordDto) (*domains.User, error) {
	var userModel User
	if err := r.db.WithContext(ctx).Joins("Country").First(&userModel, id).Error; err != nil {
		return nil, err
	}

	if dto.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		userModel.Password = string(hashedPassword)
	}

	if err := r.db.WithContext(ctx).Save(&userModel).Error; err != nil {
		return nil, err
	}

	return userModel.ToDomainWithCountryAndTimestamps(), nil
}

func (r *repository) SoftRemove(ctx context.Context, id int64) error {
	var userModel User
	if err := r.db.WithContext(ctx).First(&userModel, id).Error; err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Delete(&userModel).Error; err != nil {
		return err
	}

	return nil
}
