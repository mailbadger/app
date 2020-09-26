package params

// ForgotPassword represents request body for POST /api/forgot-password
type ForgotPassword struct {
	Email string `form:"email" validate:"required,email"`
}

func (p *ForgotPassword) TrimSpaces() {
	//no -op
}
