package params

import "strings"

// PostAuthenticate represents request body for POST /api/authenticate
type PostAuthenticate struct {
	Username string `form:"username" validate:"required,email"`
	Password string `form:"password" validate:"required"`
}

func (p *PostAuthenticate) TrimSpaces() {
	p.Password = strings.TrimSpace(p.Password)
}

// PostSignUp represents request body for POST /api/signup
type PostSignUp struct {
	Email         string `form:"email" validate:"required,email"`
	Password      string `form:"password" validate:"required,min=8"`
	TokenResponse string `form:"token_response" validate:"optional"`
}

func (p *PostSignUp) TrimSpaces() {
	p.Password = strings.TrimSpace(p.Password)
}

