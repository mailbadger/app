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

// SegmentSubs represents request body for PUT /api/segments/{id}/subscribers
type SegmentSubs struct {
	Ids []int64 `form:"ids[]" validate:"gt=0,dive,required"`
}

func (p *SegmentSubs) TrimSpaces() {
	// no op
}
