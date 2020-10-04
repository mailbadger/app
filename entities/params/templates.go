package params

import "strings"

// PostTemplate represents request body for POST /api/templates
type PostTemplate struct {
	Name    string `form:"name" validate:"required,max=191"`
	Content string `form:"content" validate:"required,html"`
	Subject string `form:"subject" validate:"required,max=191"`
}

func (p *PostTemplate) TrimSpaces() {
	p.Name = strings.TrimSpace(p.Name)
	p.Subject = strings.TrimSpace(p.Subject)
}

// PutTemplate represents request body for PUT /api/templates
type PutTemplate struct {
	Content string `json:"content" form:"content" validate:"required,html"`
	Subject string `json:"subject" form:"subject" validate:"required,max=191"`
}

func (p *PutTemplate) TrimSpaces() {
	p.Subject = strings.TrimSpace(p.Subject)
}
