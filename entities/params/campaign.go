package params

import "strings"

// Campaign represents request body for POST /api/campaigns
type Campaign struct {
	Name         string `form:"name" validate:"required,max=191"`
	TemplateName string `form:"template_name" validate:"required,max=191"`
}

func (p *Campaign) TrimSpaces() {
	p.Name = strings.TrimSpace(p.Name)
	p.TemplateName = strings.TrimSpace(p.TemplateName)
}
