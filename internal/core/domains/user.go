package domains

import (
	"golang.org/x/crypto/bcrypt"

	pkg_order "github.com/afikrim/go-hexa-template/pkg/order"
	pkg_pagination "github.com/afikrim/go-hexa-template/pkg/pagination"
)

type User struct {
	ID        uint64   `json:"id"`
	Username  string   `json:"username"`
	Phone     string   `json:"phone"`
	Email     string   `json:"email"`
	Password  string   `json:"password"`
	Fullname  string   `json:"fullname"`
	Gender    bool     `json:"gender"`
	BirthDate string   `json:"birthdate"`
	Verified  bool     `json:"verified"`
	Country   *Country `json:"country"`
	Following int64    `json:"following"`
	Followes  int64    `json:"followers"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

type UserSummary struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Fullname string `json:"fullname"`
}

type UpdateUserDto struct {
	Fullname  string `json:"fullname"`
	Gender    *bool  `json:"gender"`
	BirthDate string `json:"birthdate" validate:"regexp=^.*(?=.{8,})(?=.*[a-zA-Z])(?=.*\\d).*$"`
	CountryID uint64 `json:"country_id"`
}

type UpdateUserCredentialDto struct {
	Username string `json:"username" `
	Phone    string `json:"phone" `
	Email    string `json:"email" validate:"email"`
}

type UpdateUserPasswordDto struct {
	Password string `json:"password" validate:"regexp=^.*(?=.{8,})(?=.*[a-zA-Z])(?=.*\\d).*$"`
}

type QueryParamUserDto struct {
	Search string
	pkg_pagination.QueryParamPaginationDto
	pkg_order.QueryParamOrderDto
}

type QueryParamFollowDto struct {
	pkg_pagination.QueryParamPaginationDto
}

func (u *User) IsPasswordValid(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
}

func (u *User) HidePassword() *User {
	u.Password = ""
	return u
}
