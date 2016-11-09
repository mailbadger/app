package storage

import (
	"github.com/FilipNikolovski/news-maily/entities"
	"github.com/FilipNikolovski/news-maily/utils/pagination"
)

// GetSubscribers fetches subscribers by user id, and populates the pagination obj
func (db *store) GetSubscribers(userID int64, p *pagination.Pagination) {
	var subs []entities.Subscriber
	var count uint64

	db.Offset(p.Offset).Limit(p.PerPage).Where("user_id = ?", userID).Find(&subs).Count(&count)
	p.SetTotal(count)

	for _, t := range subs {
		p.Append(t)
	}
}

// GetSubscriber returns the subscriber by the given id and user id
func (db *store) GetSubscriber(id, userID int64) (*entities.Subscriber, error) {
	var s = new(entities.Subscriber)
	err := db.Where("user_id = ? and id = ?", userID, id).Preload("Metadata").Find(s).Error
	return s, err
}

// GetSubscriberByEmail returns the subscriber by the given email and user id
func (db *store) GetSubscriberByEmail(email string, userID int64) (*entities.Subscriber, error) {
	var s = new(entities.Subscriber)
	err := db.Where("user_id = ? and email = ?", userID, email).Preload("Metadata").Find(s).Error
	return s, err
}

func (db *store) GetSubscribersByListId(listID, userID int64, p *pagination.Pagination) {
	var l = new(entities.List)
	var count uint64

	db.Offset(p.Offset).Limit(p.PerPage).Where("user_id = ? and id = ?", userID, listID).Preload("Subscribers").Count(&count).Find(&l)
	p.SetTotal(count)

	for _, t := range l.Subscribers {
		p.Append(t)
	}
}

// CreateSubscriber creates a new subscriber in the database.
func (db *store) CreateSubscriber(s *entities.Subscriber) error {
	return db.Create(s).Error
}

// UpdateSubscriber edits an existing subscriber in the database.
func (db *store) UpdateSubscriber(s *entities.Subscriber) error {
	return db.Where("id = ? and user_id = ?", s.Id, s.UserId).Save(s).Error
}

// DeleteSubscriber deletes an existing subscriber from the database.
func (db *store) DeleteSubscriber(id, userID int64) error {
	return db.Where("user_id = ?", userID).Delete(entities.Subscriber{Id: id}).Error
}
