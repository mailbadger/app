package storage

import "github.com/news-maily/api/entities"

//CreateUser creates a new user
func (db *store) CreateUser(user *entities.User) error {
	return db.Create(user).Error
}

//UpdateUser updates the given user
func (db *store) UpdateUser(user *entities.User) error {
	return db.Save(user).Error
}

//GetUser returns an active user by id. If no user is found, an error is returned
func (db *store) GetUser(id int64) (*entities.User, error) {
	var user = new(entities.User)
	err := db.Where("id = ? and active = ?", id, true).First(user).Error
	return user, err
}

//GetUserByUsername returns a user by username. If no user is found,
//an error is returned
func (db *store) GetUserByUsername(username string) (*entities.User, error) {
	var user = new(entities.User)
	err := db.Where("username = ?", username).First(user).Error
	return user, err
}

//GetActiveUserByUsername returns an active user by username. If no user is found,
//an error is returned
func (db *store) GetActiveUserByUsername(username string) (*entities.User, error) {
	var user = new(entities.User)
	err := db.Where("username = ? and active = ?", username, true).First(user).Error
	return user, err
}
