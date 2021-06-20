package params

import (
	"strings"
)

// PostCampaign represents request body for POST /api/campaigns
type PostCampaign struct {
	Name         string `form:"name" validate:"required,max=191"`
	TemplateName string `form:"template_name" validate:"required,max=191"`
}

func (p *PostCampaign) TrimSpaces() {
	p.Name = strings.TrimSpace(p.Name)
	p.TemplateName = strings.TrimSpace(p.TemplateName)
}

// PutCampaign represents request body for PUT /api/campaigns/{id}
type PutCampaign struct {
	Name         string `form:"name" validate:"required,max=191"`
	TemplateName string `form:"template_name" validate:"required,max=191"`
}

func (p *PutCampaign) TrimSpaces() {
	p.Name = strings.TrimSpace(p.Name)
	p.TemplateName = strings.TrimSpace(p.TemplateName)
}

// StartCampaign represents request body for POST /api/campaigns/id/start
type StartCampaign struct {
	SegmentIDs          []int64           `form:"segment_id[]" validate:"required,gt=0,dive,required"`
	Source              string            `form:"source" validate:"required,email,max=191"`
	FromName            string            `form:"from_name" validate:"required,max=191"`
	DefaultTemplateData map[string]string `form:"default_template_data" validate:"dive,keys,required,alphanumhyphen,endkeys,required"`
}

func (p *StartCampaign) TrimSpaces() {
	p.FromName = strings.TrimSpace(p.FromName)
}

type CampaignSchedule struct {
	ScheduledAt         string            `form:"scheduled_at" validate:"required,datetime=2006-01-02 15:04:05,max=191"`
	FromName            string            `form:"from_name" validate:"required,max=191"`
	DefaultTemplateData map[string]string `form:"default_template_data" validate:"dive,keys,required,alphanumhyphen,endkeys,required"`
	Source              string            `form:"source" validate:"required,email,max=191"`
	SegmentIDs          []int64           `form:"segment_id[]" validate:"required,gt=0,dive,required"`
}

func (p *CampaignSchedule) TrimSpaces() {
}
