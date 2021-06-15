package entities

import (
	"encoding/json"
	"time"

	"github.com/segmentio/ksuid"
)

type CampaignSchedule struct {
	ID                      ksuid.KSUID       `json:"id" gorm:"column:id; primary_key:yes"`
	UserID                  int64             `json:"-"`
	CampaignID              int64             `json:"-"`
	ScheduledAt             time.Time         `json:"scheduled_at"`
	Source                  string            `json:"-"`
	FromName                string            `json:"-"`
	SegmentIDsJSON          JSON              `json:"-" gorm:"column:segment_ids; type:json"`
	SegmentIDs              []int64           `json:"-" sql:"-"`
	DefaultTemplateDataJSON JSON              `json:"-"  gorm:"column:default_template_data; type:json"`
	DefaultTemplateData     map[string]string `json:"-" sql:"-"`
	CreatedAt               time.Time         `json:"created_at"`
	UpdatedAt               time.Time         `json:"updated_at"`
}

// GetMetadata returns the template metadata fields.
func (s *CampaignSchedule) GetMetadata() (map[string]string, error) {
	m := make(map[string]string)

	if !s.DefaultTemplateDataJSON.IsNull() {
		err := json.Unmarshal(s.DefaultTemplateDataJSON, &m)
		if err != nil {
			return nil, err
		}
	}
	s.DefaultTemplateData = m

	return m, nil
}

// GetSegmentIDs returns the segment ids for campaign.
func (s *CampaignSchedule) GetSegmentIDs() ([]int64, error) {
	var seg []int64

	if !s.SegmentIDsJSON.IsNull() {
		err := json.Unmarshal(s.SegmentIDsJSON, &seg)
		if err != nil {
			return nil, err
		}
	}
	s.SegmentIDs = seg

	return seg, nil
}
