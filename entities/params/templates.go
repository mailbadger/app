package params

import "strings"

// PostTemplate represents request body for POST /api/templates
type PostTemplate struct {
	Name        string `json:"name" validate:"required,max=191"`
	HTMLPart    string `json:"html_part" validate:"required,html"`
	TextPart    string `json:"text_part" validate:"required"`
	SubjectPart string `json:"subject_part" validate:"required,max=191"`
}

func (p *PostTemplate) TrimSpaces() {
	p.Name = strings.TrimSpace(p.Name)
	p.SubjectPart = strings.TrimSpace(p.SubjectPart)
}

// PutTemplate represents request body for PUT /api/templates
type PutTemplate struct {
	HTMLPart    string `json:"html_part" validate:"required,html"`
	TextPart    string `json:"text_part" validate:"required"`
	SubjectPart string `json:"subject_part" validate:"required,max=191"`
	Name        string `json:"name" validate:"required,max=191"`
}

func (p *PutTemplate) TrimSpaces() {
	p.Name = strings.TrimSpace(p.Name)
	p.SubjectPart = strings.TrimSpace(p.SubjectPart)
}
