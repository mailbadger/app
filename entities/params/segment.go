package params

import (
	"strings"
)

// Segment represents request body for POST /api/segments & PUT /api/segments/{id}
type Segment struct {
	Name string `form:"name" validate:"required,max=191"`
}

func (p *Segment) TrimSpaces() {
	p.Name = strings.TrimSpace(p.Name)
}
