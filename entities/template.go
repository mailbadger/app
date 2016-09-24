package entities

import (
	"bytes"
	"errors"
	"html/template"
	"time"
)

// Template represents the template entity i.e. the email template to be sent
// to subscribers
type Template struct {
	Id        int64     `json:"id" gorm:"column:id; primary_key:yes"`
	UserId    int64     `json:"-" gorm:"column:user_id"`
	Name      string    `json:"name"`
	Content   string    `json:"content" gorm:"column:content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var ErrNameInvalid = errors.New("The name you provided is invalid.")
var ErrContentInvalid = errors.New("The content you provided is invalid.")

// Validate template properties,
// the template should be able to execute with the given variables
func (t *Template) Validate() error {
	switch {
	case t.Name == "":
		return ErrNameInvalid
	case t.Content == "":
		return ErrContentInvalid
	}

	var buff bytes.Buffer

	td := struct {
		Tracker string
		From    string
	}{
		"<img src='http://example.com/track",
		"John Doe <foo@bar.com>",
	}

	tmpl, err := template.New("html_template").Parse(t.Content)
	if err != nil {
		return err
	}

	return tmpl.Execute(&buff, td)
}
