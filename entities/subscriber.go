package entities

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"

	valid "github.com/asaskevich/govalidator"
	"github.com/mailbadger/app/utils"
)

// Subscriber represents the subscriber entity
type Subscriber struct {
	Model
	UserID      int64             `json:"-" gorm:"column:user_id; index"`
	Name        string            `json:"name"`
	Email       string            `json:"email" gorm:"not null"`
	MetaJSON    JSON              `json:"metadata" gorm:"column:metadata; type:json"`
	Segments    []Segment         `json:"segments" gorm:"many2many:subscribers_segments;"`
	Blacklisted bool              `json:"blacklisted"`
	Active      bool              `json:"active"`
	Errors      map[string]string `json:"-" sql:"-"`
	Metadata    map[string]string `json:"-" sql:"-"`
}

// GetMetadata returns the subscriber's metadata fields.
func (s *Subscriber) GetMetadata() (map[string]string, error) {
	m := make(map[string]string)

	if !s.MetaJSON.IsNull() {
		err := json.Unmarshal(s.MetaJSON, &m)
		if err != nil {
			return nil, err
		}
	}
	s.Metadata = m

	return m, nil
}

// GetUnsubscribeURL generates and signs a token based on the subscriber ID
// and creates an unsubscribe url with the email and token as query parameters.
func (s *Subscriber) GetUnsubscribeURL(uuid string) (string, error) {
	t, err := s.GenerateUnsubscribeToken(os.Getenv("UNSUBSCRIBE_SECRET"))
	if err != nil {
		return "", err
	}
	params := url.Values{}
	params.Add("email", s.Email)
	params.Add("uuid", uuid)
	params.Add("t", t)

	return os.Getenv("APP_URL") + "/unsubscribe.html?" + params.Encode(), nil
}

// Validate subscriber properties,
func (s *Subscriber) Validate() bool {
	s.Errors = make(map[string]string)

	if len(s.Name) > 0 { // Name is optional
		if valid.Trim(s.Name, "") == "" {
			s.Errors["name"] = "The subscriber name cannot be empty."
		}

		if !valid.StringLength(s.Name, "1", "191") {
			s.Errors["name"] = "The name needs to be shorter than 190 characters."
		}
	}

	if !valid.IsEmail(s.Email) {
		s.Errors["email"] = "The specified email is not valid."
	}

	for key := range s.Metadata {
		if !valid.Matches(key, "^[\\w-]*$") {
			s.Errors["message"] = fmt.Sprintf("The specified key %s must consist only of alphanumeric and hyphen characters", key)
			break
		}
	}

	return len(s.Errors) == 0
}

// GenerateUnsubscribeToken generates and signs a new unsubscribe token with the given secret, from the
// ID of the subscriber. When a subscriber wants to unsubscribe from future emails, we check this hash
// against a newly generated hash and compare them, if they match we unsubscribe the user.
func (s *Subscriber) GenerateUnsubscribeToken(secret string) (string, error) {
	if s.ID == 0 {
		return "", errors.New("entities: unable to generate unsubscribe token: subscriber ID is 0")
	}

	if secret == "" {
		return "", errors.New("entities: unable to generate unsubscribe token: secret is empty")
	}

	return utils.SignData(strconv.FormatInt(s.ID, 10), secret)
}

func (s Subscriber) GetID() int64 {
	return s.Model.ID
}

func (s Subscriber) GetCreatedAt() time.Time {
	return s.Model.CreatedAt
}

func (s Subscriber) GetUpdatedAt() time.Time {
	return s.Model.UpdatedAt
}
