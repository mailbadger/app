package entities

//Campaign represents the campaign entity
type Campaign struct {
	Id         int64  `json:"id"`
	UserId     int64  `json:"-"`
	Name       string `json:"name" sql:"not null"`
	TemplateId int64  `json:"-"`
}
