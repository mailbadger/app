package entities

import (
	"time"

	valid "github.com/asaskevich/govalidator"
)

type SesKeys struct {
	Id        int64     `json:"id" gorm:"column:id; primary_key:yes"`
	UserId    int64     `json:"-" gorm:"column:user_id; index"`
	AccessKey string    `json:"access_key" gorm:"not null"  valid:"alphanum,required"`
	SecretKey string    `json:"secret_key,omitempty" gorm:"not null" valid:"required"`
	Region    string    `json:"region" gorm:"not null" valid:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Errors map[string]string `json:"-" sql:"-"`
}

// Validate campaign properties,
func (s *SesKeys) Validate() bool {
	s.Errors = make(map[string]string)

	res, err := valid.ValidateStruct(s)
	if err != nil || !res {
		s.Errors["errors"] = err.Error()
	}

	return len(s.Errors) == 0
}
