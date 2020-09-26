package params

import (
	"strings"
)

// SegmentParams represents request body for POST /api/segments
type SegmentParams struct {
	Name string `form:"name" validate:"required,max=191"`
}

func (p *SegmentParams) TrimSpaces() {
	p.Name = strings.TrimSpace(p.Name)
}
