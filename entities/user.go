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
	BoundaryID int64          `json:"-"`
	Boundaries Boundaries     `json:"boundaries" gorm:"foreignKey:boundary_id"`
	Roles      []Role         `json:"roles"`
	Source     string         `json:"source,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}
