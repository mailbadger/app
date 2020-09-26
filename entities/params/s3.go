package params

import "strings"

// GetSignedURL represents request body for POST /api/s3/sign
type GetSignedURL struct {
	Filename string `form:"filename" validate:"required,max191"`
	ContentType string`form:"content_type" validate:"required,max=191"`
	Action string `form:"action" validate:"required,oneof=import export remove"`
}

func (p *GetSignedURL) TrimSpaces() {
	p.Filename = strings.TrimSpace(p.Filename)
	p.ContentType = strings.TrimSpace(p.ContentType)
	p.Action = strings.TrimSpace(p.Action)
}