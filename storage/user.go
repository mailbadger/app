package storage

import (
	"gorm.io/gorm"

	"github.com/mailbadger/app/entities"
)

// CreateUser creates a new user
func (db *store) CreateUser(user *entities.User) error {
	return db.Create(user).Error
}

// UpdateUser updates the given user
func (db *store) UpdateUser(user *entities.User) error {
	return db.Save(user).Error
}

// GetUser returns an active user by id. If no user is found, an error is returned
func (db *store) GetUser(id int64) (*entities.User, error) {
	var user = new(entities.User)
	err := db.
		Preload("Boundaries").
		Preload("Roles").
		Where("id = ? and active = ?", id, true).
		First(user).Error
	return user, err
}

// GetUserByUUID returns an user by uuid. If no user is found, an error is returned
func (db *store) GetUserByUUID(uuid string) (*entities.User, error) {
	var user = new(entities.User)
	err := db.
		Preload("Boundaries").
		Preload("Roles").
		Where("uuid = ?", uuid).
		First(user).Error
	return user, err
}

// GetUserByUsername returns a user by username. If no user is found,
// an error is returned
func (db *store) GetUserByUsername(username string) (*entities.User, error) {
	var user = new(entities.User)
	err := db.
		Preload("Boundaries").
		Preload("Roles").
		Where("username = ?", username).
		First(user).Error
	return user, err
}

// GetActiveUserByUsername returns an active user by username. If no user is found,
// an error is returned
func (db *store) GetActiveUserByUsername(username string) (*entities.User, error) {
	var user = new(entities.User)
	err := db.
		Preload("Boundaries").
		Preload("Roles").
		Where("username = ? and active = ?", username, true).
		First(user).Error
	return user, err
}

// DeleteUser deletes user by id
func (db *store) DeleteUser(user *entities.User) error {
	return db.Where("id = ?", user.ID).Delete(user).Error
}

// BelongsToUser finds a resource by the given user id.
func BelongsToUser(userID int64) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ?", userID)
	}
}
