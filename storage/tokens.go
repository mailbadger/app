package storage

import (
	"time"

	"github.com/mailbadger/app/entities"
)

// GetToken returns the one-time token by the given token string.
// If the token is expired, we return an error indicating that a token is not found.
func (db *store) GetToken(token string) (*entities.Token, error) {
	var t = new(entities.Token)
	err := db.Where("token = ? and expires_at > ?", token, time.Now().UTC()).First(t).Error
	if err != nil {
		return nil, err
	}
	return t, nil
}

// CreateToken adds new token in the database.
func (db *store) CreateToken(t *entities.Token) error {
	return db.Create(t).Error
}

// DeleteToken deletes the token by the given token.
func (db *store) DeleteToken(token string) error {
	return db.Delete(&entities.Token{}, "token = ?", token).Error
}

// DeleteAllTokensForUser deletes all tokens for user
func (db *store) DeleteAllTokensForUser(userID int64) error {
	return db.Where("user_id = ?", userID).Delete(&entities.Token{}).Error
}
