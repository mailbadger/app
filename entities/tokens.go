package entities

import "time"

// Different types to identify tokens.
const (
	UnsubscribeTokenType    = "unsubscribe"
	ForgotPasswordTokenType = "forgot_password"
	VerifyEmailTokenType    = "verify_email"
)

// Token entity represents a one-time token which a user can use
// in an "unsubscribe" or "forgot password" scenario.
type Token struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"-" gorm:"column:user_id; index"`
	Token     string    `json:"token" gorm:"not null" valid:"required,stringlength(1|191)"`
	Type      string    `json:"type"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
