package entities

//User represents the user entity
type User struct {
	Id       int64  `json:"id"`
	Username string `json:"username" sql:"not null;unique"`
	Password string `json:"-"`
	ApiKey   string `json:"api_key" sql:"not null;unique"`
}

//UpdateUser updates the given user
func UpdateUser(user *User) error {
	return db.Save(&user).Error
}

//GetUser returns user by id. If no user is found, an error is returned
func GetUser(id int64) (User, error) {
	user := User{}
	err := db.Where("id = ?", id).First(&user).Error
	return user, err
}

//GetUserByApiKey returns user by api key. If no user is found,
//an error is returned
func GetUserByApiKey(api_key string) (User, error) {
	user := User{}
	err := db.Where("api_key = ?", api_key).First(&user).Error
	return user, err
}
