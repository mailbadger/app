package params

import "strings"

// PutTemplate represents request body for PUT /api/templates
type PutTemplate struct {
	Content string `json:"content" form:"content" validate:"required,html"`
	Subject string `json:"subject" form:"subject" validate:"required,max=191"`
}

func (p *PutTemplate) TrimSpaces() {
	p.Content = strings.TrimSpace(p.Content)
	p.Subject = strings.TrimSpace(p.Subject)
}
