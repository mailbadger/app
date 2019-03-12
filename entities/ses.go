package entities

import (
	"time"

	valid "github.com/asaskevich/govalidator"
)

type SesKeys struct {
	Id        int64     `json:"id" gorm:"column:id; primary_key:yes"`
	UserId    int64     `json:"-" gorm:"column:user_id; index"`
	AccessKey string    `json:"access_key" gorm:"not null"  valid:"alphanum,required"`
	SecretKey string    `json:"-" gorm:"not null"  valid:"alphanum,required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Errors map[string]string `json:"-" sql:"-"`
}

// Validate campaign properties,
func (s *SesKeys) Validate() bool {
	s.Errors = make(map[string]string)

	if valid.Trim(s.AccessKey, "") == "" {
		s.Errors["name"] = ErrCampaignNameEmpty.Error()
	}
	if valid.Trim(s.SecretKey, "") == "" {
		s.Errors["subject"] = ErrSubjectEmpty.Error()
	}

	res, err := valid.ValidateStruct(s)
	if err != nil || !res {
		s.Errors["reason"] = err.Error()
	}

	return len(s.Errors) == 0
}
