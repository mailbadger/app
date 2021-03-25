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
	return db.Delete(&entities.Session{SessionID: sessionID}).Error
}
