package entities

import (
	"errors"
	"time"

	valid "github.com/asaskevich/govalidator"
)

type SesKeys struct {
	Id        int64     `json:"id" gorm:"column:id; primary_key:yes"`
	UserId    int64     `json:"-" gorm:"column:user_id; index"`
	AccessKey string    `json:"access_key" gorm:"not null"  valid:"alphanum,required"`
	SecretKey string    `json:"secret_key" gorm:"not null" valid:"alphanum,required"`
	Region    string    `json:"region" gorm:"not null" valid:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Errors map[string]string `json:"-" sql:"-"`
}

var ErrAccessKeyEmpty = errors.New("the access key cannot be empty")
var ErrSecretKeyEmpty = errors.New("the secret key cannot be empty")
var ErrRegionEmpty = errors.New("the region cannot be empty")

// Validate campaign properties,
func (s *SesKeys) Validate() bool {
	s.Errors = make(map[string]string)

	if valid.Trim(s.AccessKey, "") == "" {
		s.Errors["access_key"] = ErrAccessKeyEmpty.Error()
	}

	if valid.Trim(s.SecretKey, "") == "" {
		s.Errors["secret_key"] = ErrSecretKeyEmpty.Error()
	}

	if valid.Trim(s.SecretKey, "") == "" {
		s.Errors["region"] = ErrRegionEmpty.Error()
	}

	res, err := valid.ValidateStruct(s)
	if err != nil || !res {
		s.Errors["reason"] = err.Error()
	}

	return len(s.Errors) == 0
}
