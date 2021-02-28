package entities

type Role struct {
	ID          int64        `json:"id" gorm:"column:id; primary_key:yes"`
	Name        string       `json:"name"`
	Permissions []Permission `json:"permissions"`
}

type Permission struct {
	ID        int64  `json:"id" gorm:"column:id; primary_key:yes"`
	Name      string `json:"name"`
	Object    string `json:"object"`
	Operation string `json:"operation"`
}
