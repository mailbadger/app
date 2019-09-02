package entities

import "time"

// Session represents a user session which maps the session id stored in the cookie
// to the user that is currently signed in.
type Session struct {
	ID        int64 `gorm:"column:id; primary_key:yes"`
	UserID    int64 `gorm:"column:user_id; index"`
	User      User
	SessionID string
	CreatedAt time.Time
	UpdatedAt time.Time
}
