package entities

//User represents the user entity
type User struct {
	Id       int64  `json:"id" gorm:"column:id; primary_key:yes"`
	Username string `json:"username" gorm:"not null;unique"`
	Password string `json:"-"`
	ApiKey   string `json:"api_key" gorm:"not null;unique"`
	AuthKey  string `json:"-" gorm:"not null;unique"`
}
