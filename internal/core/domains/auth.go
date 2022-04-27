package domains

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"time"

	"github.com/golang-jwt/jwt"
)

type JwtCustomClaims struct {
	Session
	jwt.StandardClaims
}

type Session struct {
	ID           uint64 `json:"id"`
	UserID       uint64 `json:"user_id"`
	UserUsername string `json:"user_username"`
	UserEmail    string `json:"user_email"`
	UserPhone    string `json:"user_phone"`
}

type AuthWithRefresh struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IssuedAt     int64  `json:"issued_at"`
	ExpiresIn    int64  `json:"expires_in"`
}

type AuthWithoutRefresh struct {
	AccessToken string `json:"access_token"`
	IssuedAt    int64  `json:"issued_at"`
	ExpiresIn   int64  `json:"expires_in"`
}

type RegisterDto struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Password  string `json:"password"`
	Fullname  string `json:"fullname"`
	Gender    bool   `json:"gender"`
	BirthDate string `json:"birthdate"`
	CountryID uint64 `json:"country_id"`
}

type LoginDto struct {
	Credential string `json:"credential"`
	Password   string `json:"password"`
}

func (s *Session) GenerateRefreshToken(secret string) string {
	queryVal := url.Values{}
	queryVal.Add("username", s.UserUsername)
	queryVal.Add("email", s.UserEmail)
	queryVal.Add("phone", s.UserPhone)

	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(queryVal.Encode()))

	return hex.EncodeToString(h.Sum(nil))
}

func (s *Session) GenerateAccessToken(secret string, expiresIn int64) (*string, error) {
	claims := &JwtCustomClaims{
		Session{
			ID:           s.ID,
			UserID:       s.UserID,
			UserUsername: s.UserUsername,
			UserEmail:    s.UserEmail,
			UserPhone:    s.UserPhone,
		},
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * time.Duration(expiresIn)).Unix(),
		},
	}

	newJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := newJwt.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}

	return &token, nil
}
