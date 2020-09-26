package params

import "strings"

// SendCampaign represents request body for POST /api/campaigns/id/start
type SendCampaign struct {
	Ids      []int64 `form:"segment_id[]" validate:"gt=0,dive,required"`
	Source   string  `form:"source" validate:"required,email,max=191"`
	FromName string  `form:"from_name" validate:"required,max=191"`
}

func (p *SendCampaign) TrimSpaces() {
	p.FromName = strings.TrimSpace(p.FromName)
}
