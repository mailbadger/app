package params

// PutForgotPassword represents request body for PUT /api/forgot-password/{token}
type PutForgotPassword struct {
	Password string `form:"password" validate:"required,min=8"`
}

func (p *PutForgotPassword) TrimSpaces() {
	// no-op
}
