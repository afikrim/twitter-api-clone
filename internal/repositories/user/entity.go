package user_repository

import (
	"time"

	"github.com/afikrim/go-hexa-template/internal/core/domains"
	country_repository "github.com/afikrim/go-hexa-template/internal/repositories/country"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID             uint64                     `gorm:"column:id;not null;primaryKey;autoIncrement"`
	Username       string                     `gorm:"column:username;type:varchar(255);not null;unique"`
	Phone          string                     `gorm:"column:phone;type:varchar(255);not null;unique"`
	Email          string                     `gorm:"column:email;type:varchar(255);not null;unique"`
	Password       string                     `gorm:"column:password;type:varchar(255);not null"`
	Fullname       string                     `gorm:"column:fullname;type:varchar(255);not null"`
	Gender         bool                       `gorm:"column:gender;not null;default:false"`
	Verified       bool                       `gorm:"column:verified;not null;default:false"`
	BirthDate      time.Time                  `gorm:"column:birthdate;type:date;not null"`
	CountryID      uint64                     `gorm:"column:country_id;type:bigint;not null"`
	Country        country_repository.Country `gorm:"foreignKey:country_id;references:id;target:countries"`
	Following      []User                     `gorm:"many2many:user_following;foreignKey:id;joinForeignKey:following_id;references:id;joinReferences:follower_id"`
	FollowingCount int64                      `gorm:"-"`
	Followers      []User                     `gorm:"many2many:user_following;foreignKey:id;joinForeignKey:following_id;references:id;joinReferences:follower_id"`
	FollowersCount int64                      `gorm:"-"`
	CreatedAt      *time.Time                 `gorm:"column:created_at;not null;autoCreateTime"`
	UpdatedAt      *time.Time                 `gorm:"column:updated_at;not null;autoUpdateTime"`
	DeletedAt      gorm.DeletedAt             `gorm:"column:deleted_at"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) ToDomain() *domains.User {
	return &domains.User{
		ID:        u.ID,
		Username:  u.Username,
		Phone:     u.Phone,
		Email:     u.Email,
		Password:  u.Password,
		Fullname:  u.Fullname,
		Gender:    u.Gender,
		Verified:  u.Verified,
		BirthDate: u.BirthDate.Format("2006-01-02"),
	}
}

func (u *User) ToDomainSummary() *domains.UserSummary {
	return &domains.UserSummary{
		ID:       u.ID,
		Username: u.Username,
		Fullname: u.Fullname,
	}
}

func (u *User) ToDomainWithTimestamps() *domains.User {
	user := u.ToDomain()
	user.CreatedAt = u.CreatedAt.Format("2006-01-02 15:04:05")
	user.UpdatedAt = u.UpdatedAt.Format("2006-01-02 15:04:05")

	return user
}

func (u *User) ToDomainWithCountry() *domains.User {
	user := u.ToDomain()
	user.Country = u.Country.ToDomain()

	return user
}

func (u *User) ToDomainWithCountryAndTimestamps() *domains.User {
	user := u.ToDomainWithCountry()
	user.CreatedAt = u.CreatedAt.Format("2006-01-02 15:04:05")
	user.UpdatedAt = u.UpdatedAt.Format("2006-01-02 15:04:05")

	return user
}

func (u *User) ToDomainDetail() *domains.User {
	user := u.ToDomainWithCountryAndTimestamps()
	user.Followes = u.FollowingCount
	user.Following = u.FollowersCount

	return user
}

func (User) FromRegisterDto(d *domains.RegisterDto) *User {
	parsedBirthDate, err := time.Parse("2006-01-02", d.BirthDate)
	if err != nil {
		panic(err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(d.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	return &User{
		Username:  d.Username,
		Phone:     d.Phone,
		Email:     d.Email,
		Password:  string(hashedPassword),
		Fullname:  d.Fullname,
		Gender:    d.Gender,
		CountryID: d.CountryID,
		BirthDate: parsedBirthDate,
	}
}

func (User) FromDomain(d *domains.User) *User {
	parsedBirthDate, err := time.Parse("2006-01-02", d.BirthDate)
	if err != nil {
		panic(err)
	}

	return &User{
		ID:        d.ID,
		Username:  d.Username,
		Phone:     d.Phone,
		Email:     d.Email,
		Password:  d.Password,
		Fullname:  d.Fullname,
		Gender:    d.Gender,
		Verified:  d.Verified,
		BirthDate: parsedBirthDate,
		Country:   country_repository.Country{ID: d.Country.ID, Name: d.Country.Name, Code: d.Country.Code},
	}
}

func (User) FromDomainWithTimestamps(d *domains.User) *User {
	parsedBirthDate, err := time.Parse("2006-01-02", d.BirthDate)
	if err != nil {
		panic(err)
	}

	parsedCreatedAt, err := time.ParseInLocation("2006-01-02 15:04:05", d.CreatedAt, time.UTC)
	if err != nil {
		panic(err)
	}

	parsedUpdatedAt, err := time.ParseInLocation("2006-01-02 15:04:05", d.UpdatedAt, time.UTC)
	if err != nil {
		panic(err)
	}

	return &User{
		ID:        d.ID,
		Username:  d.Username,
		Phone:     d.Phone,
		Email:     d.Email,
		Password:  d.Password,
		Fullname:  d.Fullname,
		Gender:    d.Gender,
		Verified:  d.Verified,
		BirthDate: parsedBirthDate,
		Country:   country_repository.Country{ID: d.Country.ID, Name: d.Country.Name, Code: d.Country.Code},
		CreatedAt: &parsedCreatedAt,
		UpdatedAt: &parsedUpdatedAt,
	}
}
