package userfollowing_repository

import (
	"context"

	"github.com/afikrim/go-hexa-template/internal/core/domains"
	user_repository "github.com/afikrim/go-hexa-template/internal/repositories/user"
	pkg_pagination "github.com/afikrim/go-hexa-template/pkg/pagination"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewUserFollowingRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, currentUser uint64, followUser uint64) error {
	userFollowingModel := map[string]interface{}{
		"following_id": followUser,
		"follower_id":  currentUser,
	}
	if err := r.db.WithContext(ctx).Table("user_following").Create(&userFollowingModel).Error; err != nil {
		return err
	}

	return nil
}

func (r *repository) FindAllFollowing(ctx context.Context, userID uint64, query *domains.QueryParamFollowDto) ([]domains.UserSummary, *pkg_pagination.CursorPagination, error) {
	var userModels []user_repository.User
	var userSummaries []domains.UserSummary
	var cursorPagination *pkg_pagination.CursorPagination
	var countTotal int64

	if err := r.db.WithContext(ctx).Model(&userModels).
		Joins("RIGHT JOIN user_following ON user_following.following_id = users.id").
		Where("user_following.follower_id = ?", userID).
		Order("users.fullname asc").
		Limit(*query.Limit).
		Offset(*query.Offset).
		Find(&userModels).Error; err != nil {
		return nil, nil, err
	}

	for _, userModel := range userModels {
		userSummaries = append(userSummaries, *userModel.ToDomainSummary())
	}

	if err := r.db.WithContext(ctx).Model(&userModels).
		Joins("RIGHT JOIN user_following ON user_following.following_id = users.id").
		Where("user_following.follower_id = ?", userID).
		Count(&countTotal).Error; err != nil {
		return nil, nil, err
	}

	cursorPagination = pkg_pagination.NewCursorPagination(countTotal, *query.Limit, *query.Offset)

	return userSummaries, cursorPagination, nil
}

func (r *repository) FindAllFollowers(ctx context.Context, userID uint64, query *domains.QueryParamFollowDto) ([]domains.UserSummary, *pkg_pagination.CursorPagination, error) {
	var userModels []user_repository.User
	var userSummaries []domains.UserSummary
	var cursorPagination *pkg_pagination.CursorPagination
	var countTotal int64

	if err := r.db.WithContext(ctx).Model(&userModels).
		Joins("RIGHT JOIN user_following ON user_following.follower_id = users.id").
		Where("user_following.following_id = ?", userID).
		Order("users.fullname asc").
		Limit(*query.Limit).
		Offset(*query.Offset).
		Find(&userModels).Error; err != nil {
		return nil, nil, err
	}

	for _, userModel := range userModels {
		userSummaries = append(userSummaries, *userModel.ToDomainSummary())
	}

	if err := r.db.WithContext(ctx).Model(&userModels).
		Joins("RIGHT JOIN user_following ON user_following.follower_id = users.id").
		Where("user_following.following_id = ?", userID).
		Count(&countTotal).Error; err != nil {
		return nil, nil, err
	}

	cursorPagination = pkg_pagination.NewCursorPagination(countTotal, *query.Limit, *query.Offset)

	return userSummaries, cursorPagination, nil
}

func (r *repository) Remove(ctx context.Context, currentUserID uint64, followUserID uint64) error {
	userFollowingModel := map[string]interface{}{
		"following_id": followUserID,
		"follower_id":  currentUserID,
	}
	if err := r.db.WithContext(ctx).Table("user_following").
		Where("following_id = ? AND follower_id = ?", followUserID, currentUserID).
		Delete(&userFollowingModel).Error; err != nil {
		return err
	}

	return nil
}
