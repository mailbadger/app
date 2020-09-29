package params

import (
	"strings"
)

// PostSESKeys represents request body for POST /api/ses/keys
type PostSESKeys struct {
	AccessKey string `form:"access_key" validate:"required,alphanum,max=191"`
	SecretKey string `form:"secret_key" validate:"required,max=191"`
	Region    string `form:"region" validate:"required,max=30"`
}

func (p *PostSESKeys) TrimSpaces() {
	p.AccessKey = strings.TrimSpace(p.AccessKey)
	p.SecretKey = strings.TrimSpace(p.SecretKey)
	p.Region = strings.TrimSpace(p.Region)
}
