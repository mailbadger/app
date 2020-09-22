package params

import (
	"strings"
)

// PostSESKeys represents request body for POST /api/ses/keys
type PostSESKeys struct {
	AccessKey string `form:"access_key" validate:"required,alphanum"`
	SecretKey string `form:"secret_key" validate:"required"`
	Region    string `form:"region" validate:"required"`
}

func (p *PostSESKeys) TrimSpaces() {
	p.AccessKey = strings.TrimSpace(p.AccessKey)
	p.SecretKey = strings.TrimSpace(p.SecretKey)
	p.Region = strings.TrimSpace(p.Region)
}