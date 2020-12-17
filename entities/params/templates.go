package params

import "strings"

// PostTemplate represents request body for POST /api/templates
type PostTemplate struct {
	Name     string `form:"name" validate:"required,max=191"`
	HTMLPart string `form:"html_part" validate:"required,html"`
	TextPart string `form:"text_part" validate:"required"`
	Subject  string `form:"subject" validate:"required,max=191"`
}

func (p *PostTemplate) TrimSpaces() {
	p.Name = strings.TrimSpace(p.Name)
	p.Subject = strings.TrimSpace(p.Subject)
}

// PutTemplate represents request body for PUT /api/templates
type PutTemplate struct {
	HTMLPart string `form:"html_part" validate:"required,html"`
	TextPart string `form:"text_part" validate:"required"`
	Subject  string `json:"subject" form:"subject" validate:"required,max=191"`
}

func (p *PutTemplate) TrimSpaces() {
	p.Subject = strings.TrimSpace(p.Subject)
}
