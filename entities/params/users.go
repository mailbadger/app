package params

// ChangePassword represents request body for POST /api/users/password
type ChangePassword struct {
	Password    string `json:"password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

func (p *ChangePassword) TrimSpaces() {
	//no -op
}

// PutForgotPassword represents request body for PUT /api/forgot-password/{token}
type PutForgotPassword struct {
	Password string `json:"password" validate:"required,min=8"`
}

func (p *PutForgotPassword) TrimSpaces() {
	// no-op
}

// ForgotPassword represents request body for POST /api/forgot-password
type ForgotPassword struct {
	Email string `json:"email" validate:"required,email"`
}

func (p *ForgotPassword) TrimSpaces() {
	//no -op
}
