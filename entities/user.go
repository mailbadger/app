package entities

import (
	"database/sql"
	"time"
)

// User represents the user entity
type User struct {
	ID         int64          `json:"-" gorm:"column:id; primary_key:yes"`
	UUID       string         `json:"uuid"`
	Username   string         `json:"username" gorm:"not null;unique"`
	Password   sql.NullString `json:"-"`
	Active     bool           `json:"active"`
	Verified   bool           `json:"verified"`
	Boundaries *Boundaries    `json:"-" gorm:"foreignKey:boundary_id"`
	Source     string         `json:"source,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}
