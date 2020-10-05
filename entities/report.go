package entities

const (
	StatusFailed     = "failed"
	StatusDone       = "done"
	StatusInProgress = "in_progress"

	SubscribersResource = "subscribers"
)

//Report represents the Report entity
type Report struct {
	Model
	UserID   int64  `json:"-" gorm:"column:user_id; index"`
	Resource string `json:"resource" gorm:"not null"`
	FileName string `json:"file_name" gorm:"not null"`
	Type     string `json:"type" gorm:"not null"`
	Status   string `json:"status" gorm:"not null"`
	Note     string `json:"note"`
}
