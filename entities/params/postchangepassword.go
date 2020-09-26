package params

// ChangePassword represents request body for POST /api/users/password
type ChangePassword struct {
	Password    string `form:"password" validate:"required"`
	NewPassword string `form:"new_password" validate:"required,min=8"`
}

func (p *ChangePassword) TrimSpaces() {
	//no -op
}
