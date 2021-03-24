package entities

type Role struct {
	ID   int64  `json:"id" gorm:"column:id; primary_key:yes"`
	Name string `json:"name"`
}
