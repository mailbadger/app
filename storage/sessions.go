package storage

import (
	"github.com/mailbadger/app/entities"
)

// GetSession returns the session by the given session id.
func (db *store) GetSession(sessionID string) (*entities.Session, error) {
	var s = new(entities.Session)
	err := db.Where("session_id = ?", sessionID).
		Preload("User.Boundaries").Preload("User.Roles").
		First(s).
		Error
	if err != nil {
		return nil, err
	}
	return s, nil
}

// CreateSession adds a new session in the database.
func (db *store) CreateSession(s *entities.Session) error {
	return db.Create(s).Error
}

// DeleteSession deletes a session by the given session id from the database.
func (db *store) DeleteSession(sessionID string) error {
	return db.Where("session_id = ?", sessionID).Delete(&entities.Session{}).Error
}

// DeleteAllSessionsForUser deletes sessions for user
func (db *store) DeleteAllSessionsForUser(userID int64) error {
	return db.Where("user_id = ?", userID).Delete(&entities.Session{}).Error
}
