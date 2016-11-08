package entities

import (
	"errors"
	"time"

	valid "github.com/asaskevich/govalidator"
)

// Template represents the template entity i.e. the email template to be sent
// to subscribers
type Template struct {
	Id        int64             `json:"id" gorm:"column:id; primary_key:yes"`
	UserId    int64             `json:"-" gorm:"column:user_id; index"`
	Name      string            `json:"name"`
	Content   string            `json:"content" gorm:"column:content"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	Errors    map[string]string `json:"-" sql:"-"`
}

var ErrTemplateNameEmpty = errors.New("The name cannot be empty.")
var ErrContentEmpty = errors.New("The content cannot be empty.")

// Validate template properties,
// the template should be able to execute with the given variables
func (t *Template) Validate() bool {
	t.Errors = make(map[string]string)

	if valid.Trim(t.Name, "") == "" {
		t.Errors["name"] = ErrTemplateNameEmpty.Error()
	}
	if valid.Trim(t.Content, "") == "" {
		t.Errors["content"] = ErrContentEmpty.Error()
	}

	return len(t.Errors) == 0
}
