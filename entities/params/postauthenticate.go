package params

// PostAuthenticate represents request body for POST /api/authenticate
type PostAuthenticate struct {
	Username string `form:"username" validate:"required"`
	Password string `form:"password" validate:"required"`
}

func (p *PostAuthenticate) TrimSpaces() {
}