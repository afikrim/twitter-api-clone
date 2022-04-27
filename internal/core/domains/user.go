package domains

import "golang.org/x/crypto/bcrypt"

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
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

type UpdateUserDto struct {
	Fullname  string `json:"fullname"`
	Gender    *bool  `json:"gender"`
	BirthDate string `json:"birthdate"`
	CountryID uint64 `json:"country_id"`
}

type UpdateUserCredentialDto struct {
	Username string `json:"username"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
}

type UpdateUserPasswordDto struct {
	Password string `json:"password"`
}

type QueryParamDto struct {
	Search  string
	OrderBy string
	SortBy  string
	Offset  *int
	Limit   *int
	Page    *int
}

func (u *User) IsPasswordValid(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
}

func (u *User) HidePassword() *User {
	u.Password = ""
	return u
}
