package params

import "strings"

// PostTemplate represents request body for POST /api/templates
type PostTemplate struct {
	Name        string `form:"name" validate:"required,max=191"`
	HTMLPart    string `form:"html_part" validate:"required,html"`
	TextPart    string `form:"text_part" validate:"required"`
	SubjectPart string `form:"subject_part" validate:"required,max=191"`
}

func (p *PostTemplate) TrimSpaces() {
	p.Name = strings.TrimSpace(p.Name)
	p.SubjectPart = strings.TrimSpace(p.SubjectPart)
}

// PutTemplate represents request body for PUT /api/templates
type PutTemplate struct {
	HTMLPart    string `form:"html_part" validate:"required,html"`
	TextPart    string `form:"text_part" validate:"required"`
	SubjectPart string `form:"subject_part" validate:"required,max=191"`
	Name        string `form:"name" validate:"required,max=191"`
}

func (p *PutTemplate) TrimSpaces() {
	p.Name = strings.TrimSpace(p.Name)
	p.SubjectPart = strings.TrimSpace(p.SubjectPart)
}
