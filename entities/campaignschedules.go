package entities

import (
	"encoding/json"
	"time"

	"github.com/segmentio/ksuid"
)

type CampaignSchedule struct {
	ID                      ksuid.KSUID       `json:"id" gorm:"column:id; primary_key:yes"`
	UserID                  int64             `json:"user_id"`
	CampaignID              int64             `json:"campaign_id"`
	ScheduledAt             time.Time         `json:"scheduled_at"`
	Source                  string            `json:"source"`
	FromName                string            `json:"from_name"`
	SegmentIDsJSON          JSON              `json:"segment_ids"`
	SegmentIDs              []int64           `json:"-"`
	DefaultTemplateDataJSON JSON              `json:"default_template_data"`
	DefaultTemplateData     map[string]string `json:"-"`
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
