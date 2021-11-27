package params

// PostAuthenticate represents request body for POST /api/authenticate
type PostAuthenticate struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (p *PostAuthenticate) TrimSpaces() {
	// no trimming needed
}

// PostSignUp represents request body for POST /api/signup
type PostSignUp struct {
	Email         string `json:"email" validate:"required,email"`
	Password      string `json:"password" validate:"required,min=8"`
	TokenResponse string `json:"token_response" validate:"omitempty"`
}

func (p *PostSignUp) TrimSpaces() {
	// no trimming needed
}
