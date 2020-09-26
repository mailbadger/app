package params

import "strings"

// CampaignParams represents request body for POST /api/campaigns
type CampaignParams struct {
	Name         string `form:"name" validate:"required,max=191"`
	TemplateName string `form:"template_name" validate:"required,max=191"`
}

func (p *CampaignParams) TrimSpaces() {
	p.Name = strings.TrimSpace(p.Name)
	p.TemplateName = strings.TrimSpace(p.TemplateName)
}

// CampaignParams represents request body for POST /api/campaigns/id/start
type SendCampaignParams struct {
	Ids      []int64 `form:"segment_id[]" validate:"gt=0,dive,required"`
	Source   string  `form:"source" validate:"required,email,max=191"`
	FromName string  `form:"from_name" validate:"required,max=191"`
}

func (p *SendCampaignParams) TrimSpaces() {
	p.FromName = strings.TrimSpace(p.FromName)
}
