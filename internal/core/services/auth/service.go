package auth_service

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/afikrim/go-hexa-template/internal/core/domains"
	"github.com/afikrim/go-hexa-template/internal/core/ports/repositories"
)

const (
	accessTokenExpiresIn = int64(60 * 60 * 24 * 7)
)

var (
	ErrUserNotVerified = errors.New("user not verified")
	ErrInvalidPassword = errors.New("invalid password")
	ErrSessionNotFound = errors.New("session not found")
)

type service struct {
	userRepo    repositories.UserRepository
	sessionRepo repositories.SessionRepository
}

func NewAuthService(userRepo repositories.UserRepository, sessionRepo repositories.SessionRepository) *service {
	return &service{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

func (s *service) Register(ctx context.Context, dto *domains.RegisterDto) error {
	_, err := s.userRepo.Create(ctx, dto)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) Login(ctx context.Context, dto *domains.LoginDto) (*domains.AuthWithRefresh, error) {
	var user *domains.User
	var err error

	phoneRegexp := regexp.MustCompile(`^\d+$`)
	emailRegexp := regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`)
	if phoneRegexp.MatchString(dto.Credential) {
		user, err = s.userRepo.FindByPhone(ctx, dto.Credential)
	}
	if user == nil && emailRegexp.MatchString(dto.Credential) {
		user, err = s.userRepo.FindByEmail(ctx, dto.Credential)
	}
	if user == nil {
		user, err = s.userRepo.FindByUsername(ctx, dto.Credential)
	}

	if err != nil {
		return nil, err
	}

	if !user.Verified {
		return nil, ErrUserNotVerified
	}

	if !user.IsPasswordValid(dto.Password) {
		return nil, ErrInvalidPassword
	}

	sessionIDRaw := fmt.Sprint(time.Now().Second()) + fmt.Sprint(user.ID)
	sessionID, err := strconv.ParseUint(sessionIDRaw, 10, 64)
	if err != nil {
		return nil, err
	}

	session := &domains.Session{
		ID:           sessionID,
		UserID:       user.ID,
		UserUsername: user.Username,
		UserPhone:    user.Phone,
		UserEmail:    user.Email,
	}
	refreshToken := session.GenerateRefreshToken("secret")

	err = s.sessionRepo.Create(ctx, refreshToken, session)
	if err != nil {
		return nil, err
	}

	accessToken, err := session.GenerateAccessToken("secret", accessTokenExpiresIn)
	if err != nil {
		return nil, err
	}

	return &domains.AuthWithRefresh{
		AccessToken:  *accessToken,
		RefreshToken: refreshToken,
		IssuedAt:     time.Now().Unix(),
		ExpiresIn:    accessTokenExpiresIn,
	}, nil
}

func (s *service) Refresh(ctx context.Context, refreshToken string) (*domains.AuthWithoutRefresh, error) {
	session, err := s.sessionRepo.FindByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, ErrSessionNotFound
	}

	accessToken, err := session.GenerateAccessToken("secret", accessTokenExpiresIn)
	if err != nil {
		return nil, err
	}

	return &domains.AuthWithoutRefresh{
		AccessToken: *accessToken,
		IssuedAt:    time.Now().Unix(),
		ExpiresIn:   accessTokenExpiresIn,
	}, nil
}

func (s *service) Logout(ctx context.Context, refreshToken string) error {
	session, err := s.sessionRepo.FindByRefreshToken(ctx, refreshToken)
	if err != nil {
		return err
	}
	if session == nil {
		return ErrSessionNotFound
	}

	return s.sessionRepo.Remove(ctx, refreshToken)
}
