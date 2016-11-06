package entities

import "time"

//Subscriber represents the subscriber entity
type Subscriber struct {
	Id        int64                `json:"id" gorm:"column:id; primary_key:yes"`
	Name      string               `json:"name" gorm:"not null"`
	Email     string               `json:"email" gorm:"not null"`
	Lists     []List               `gorm:"many2many:subscribers_lists;"`
	Metadata  []SubscriberMetadata `json:"metadata" gorm:"ForeignKey:SubscriberId"`
	CreatedAt time.Time            `json:"created_at"`
	UpdatedAt time.Time            `json:"updated_at"`
	Errors    map[string]string    `json:"-" sql:"-"`
}

//SubscriberMetadata represents the subscriber metadata in a form of a key and value
type SubscriberMetadata struct {
	Id           int64  `gorm:"column:id; primary_key:yes"`
	SubscriberId int64  `gorm:"column:subscriber_id; index"`
	Key          string `gorm:"not null"`
	Value        string `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
