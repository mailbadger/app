package entities

import (
	"time"

	valid "github.com/asaskevich/govalidator"
)

// Segment represents the list entity
type Segment struct {
	Model
	Name        string            `json:"name" gorm:"not null" valid:"required,stringlength(1|191)"`
	UserID      int64             `json:"-" gorm:"column:user_id; index"`
	Subscribers []Subscriber      `json:"-" gorm:"many2many:subscribers_segments;"`
	Errors      map[string]string `json:"-" sql:"-"`
}

// SegmentWithTotalSubs represents the segment entity with
// extra information regarding the total count of subscribers.
// this entity is needed because we run a custom query for the paginated
// set of results, which differs from the rest of the CRUD methods where
// 'subscribers_in_segment' column is not present.
type SegmentWithTotalSubs struct {
	Segment
	SubscribersInSeg int64  `json:"subscribers_in_segment" gorm:"column:subscribers_in_segment"`
	TotalSubscribers *int64 `json:"total_subscribers,omitempty" sql:"-"`
}

// Validate validates the list properties and populates the Errors map
// in case of any errors.
func (l *Segment) Validate() bool {
	l.Errors = make(map[string]string)

	if valid.Trim(l.Name, "") == "" {
		l.Errors["name"] = "The segment name cannot be empty."
	}

	res, err := valid.ValidateStruct(l)
	if err != nil || !res {
		l.Errors["message"] = err.Error()
	}

	return len(l.Errors) == 0
}

func (s Segment) GetID() int64 {
	return s.Model.ID
}

func (s Segment) GetCreatedAt() time.Time {
	return s.Model.CreatedAt
}

func (s Segment) GetUpdatedAt() time.Time {
	return s.Model.UpdatedAt
}
