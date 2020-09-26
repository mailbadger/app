package params

// SegmentSubs represents request body for PUT /api/segments/{id}/subscribers
type SegmentSubs struct {
	Ids []int64 `form:"ids[]" validate:"gt=0,dive,required"`
}

func (p *SegmentSubs) TrimSpaces() {
	// no op
}
