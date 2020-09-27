package params

import "strings"

// PostSubscriber represents request body for POST /api/subscribers
type PostSubscriber struct {
	Name       string            `form:"name" validate:"omitempty,min=1,max=191"`
	Email      string            `form:"email" validate:"required,email"`
	SegmentIDs []int64           `form:"segments[]" validate:"omitempty"`
	Metadata   map[string]string `form:"metadata" validate:"omitempty,dive,alphanumhyphen"`
}

func (p *PostSubscriber) TrimSpaces() {
	p.Name = strings.TrimSpace(p.Name)
}
