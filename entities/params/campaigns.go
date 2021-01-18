package params

import "strings"

// Campaign represents request body for POST /api/campaigns
type Campaign struct {
	Name         string `form:"name" validate:"required,max=191"`
	TemplateName string `form:"template_name" validate:"required,max=191"`
}

func (p *Campaign) TrimSpaces() {
	p.Name = strings.TrimSpace(p.Name)
	p.TemplateName = strings.TrimSpace(p.TemplateName)
}

// SendCampaign represents request body for POST /api/campaigns/id/start
type SendCampaign struct {
	SegmentIDs          []int64           `form:"segment_id[]" validate:"required,gt=0,dive,required"`
	Source              string            `form:"source" validate:"required,email,max=191"`
	FromName            string            `form:"from_name" validate:"required,max=191"`
	DefaultTemplateData map[string]string `form:"default_template_data" validate:"dive,keys,required,alphanumhyphen,endkeys,required"`
}

func (p *SendCampaign) TrimSpaces() {
	p.FromName = strings.TrimSpace(p.FromName)
}
