package storage

import "github.com/FilipNikolovski/news-maily/entities"

//UpdateUser updates the given user
func (db *store) UpdateUser(user *entities.User) error {
	return db.Save(user).Error
}

//GetUser returns user by id. If no user is found, an error is returned
func (db *store) GetUser(id int64) (*entities.User, error) {
	var user = new(entities.User)
	err := db.Where("id = ?", id).First(user).Error
	return user, err
}

//GetUserByAPIKey returns user by api key. If no user is found,
//an error is returned
func (db *store) GetUserByAPIKey(apiKey string) (*entities.User, error) {
	var user = new(entities.User)
	err := db.Where("api_key = ?", apiKey).First(user).Error
	return user, err
}

//GetUserByUsername returns user by username. If no user is found,
//an error is returned
func (db *store) GetUserByUsername(username string) (*entities.User, error) {
	var user = new(entities.User)
	err := db.Where("username = ?", username).First(user).Error
	return user, err
}
