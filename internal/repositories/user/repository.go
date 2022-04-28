package user_repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/afikrim/go-hexa-template/internal/core/domains"
	pkg_pagination "github.com/afikrim/go-hexa-template/pkg/pagination"
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

func (r *repository) FindAll(ctx context.Context, query *domains.QueryParamUserDto) ([]domains.UserSummary, *pkg_pagination.CursorPagination, error) {
	qb := r.db.Model(&User{}).WithContext(ctx)
	countQb := qb
	if query.Search != "" {
		search := fmt.Sprintf("%%%s%%", query.Search)
		qb = qb.Where("username LIKE ? OR fullname LIKE ?", search, search)
		countQb = qb
	}

	orderby := fmt.Sprintf("%s %s", query.SortBy, query.OrderBy)
	qb.Limit(*query.Limit)
	qb.Offset(*query.Offset)
	qb.Order(orderby)

	var userModels []User
	var count int64
	err, countErr := qb.Find(&userModels).Error, countQb.Error
	if err != nil {
		return nil, nil, err
	}
	if countErr != nil {
		return nil, nil, countErr
	}

	var users []domains.UserSummary
	for _, userModel := range userModels {
		users = append(users, *userModel.ToDomainSummary())
	}

	next := int64(*query.Offset) + int64(*query.Limit)
	if next >= count {
		next = -1
	}
	cursor := &pkg_pagination.CursorPagination{
		Current: int64(*query.Offset),
		Next:    next,
	}

	return users, cursor, nil
}

func (r *repository) FindByID(ctx context.Context, id int64) (*domains.User, error) {
	var userModel User
	if err := r.db.WithContext(ctx).First(&userModel, id).Error; err != nil {
		return nil, err
	}

	return userModel.ToDomainWithTimestamps(), nil
}

func (r *repository) FindByUsername(ctx context.Context, username string) (*domains.User, error) {
	qb := r.db.WithContext(ctx).Model(&User{})
	qb.Where("username = ?", username)
	qb.Joins("Country")

	var userModel User
	if err := qb.Scan(&userModel).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Table("user_following").Where("follower_id = ?", userModel.ID).Count(&userModel.FollowingCount).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Table("user_following").Where("following_id = ?", userModel.ID).Count(&userModel.FollowersCount).Error; err != nil {
		return nil, err
	}

	return userModel.ToDomainDetail(), nil
}

func (r *repository) FindByCredential(ctx context.Context, credential string) (*domains.User, error) {
	var userModel User
	if err := r.db.WithContext(ctx).Where("username = ? OR email = ? OR phone = ?", credential, credential, credential).First(&userModel).Error; err != nil {
		return nil, err
	}

	return userModel.ToDomainWithCountryAndTimestamps(), nil
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
