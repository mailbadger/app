package params

type SegmentSubs struct {
	Ids []int64 `form:"ids[]" validate:"required"`
}

func (p *SegmentSubs) TrimSpaces() {
}
