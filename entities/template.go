package entities

import "time"

type BaseTemplate struct {
	Model
	UserID      int64  `json:"user_id"`
	Name        string    `json:"name"`
	SubjectPart string    `json:"subject_part"`
}

func (c BaseTemplate) GetID() int64 {
	return c.ID
}

type Template struct {
	BaseTemplate
	HTMLPart    string `json:"html_part" gorm:"-"`
	TextPart    string `json:"text_part"`
}

func (t Template) GetBase() *BaseTemplate {
	return &BaseTemplate{
		Model:       Model{
			ID:        t.ID,
			CreatedAt: t.CreatedAt,
			UpdatedAt: t.UpdatedAt,
		},
		UserID:      t.UserID,
		Name:        t.Name,
		SubjectPart: t.SubjectPart,
	}
}

type TemplateCollection struct {
	NextToken  string         `json:"next_token"`
	Collection []TemplateMeta `json:"collection"`
}

type TemplateMeta struct {
	Name      string    `json:"name"`
	Timestamp time.Time `json:"timestamp"`
}
