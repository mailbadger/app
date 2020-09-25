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
	p.Content = strings.TrimSpace(p.Content)
	p.Subject = strings.TrimSpace(p.Subject)
}
