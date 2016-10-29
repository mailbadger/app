package entities

//User represents the user entity
type User struct {
	Id       int64  `json:"id"`
	Username string `json:"username" sql:"not null;unique"`
	Password string `json:"-"`
	ApiKey   string `json:"api_key" sql:"not null;unique"`
	AuthKey  string `json:"-" sql:"not null; unique"`
}
