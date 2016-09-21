package storage

import "github.com/FilipNikolovski/news-maily/entities"

//UpdateUser updates the given user
func (db *store) UpdateUser(user *entities.User) error {
	return db.Save(&user).Error
}

//GetUser returns user by id. If no user is found, an error is returned
func (db *store) GetUser(id int64) (entities.User, error) {
	user := entities.User{}
	err := db.Where("id = ?", id).First(&user).Error
	return user, err
}

//GetUserByApiKey returns user by api key. If no user is found,
//an error is returned
func (db *store) GetUserByApiKey(api_key string) (entities.User, error) {
	user := entities.User{}
	err := db.Where("api_key = ?", api_key).First(&user).Error
	return user, err
}
