package storage

import (
	"github.com/news-maily/api/entities"
	"github.com/news-maily/api/utils/pagination"
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

// GetSubscribersByIDs returns the subscriber by the given id and user id
func (db *store) GetSubscribersByIDs(ids []int64, userID int64) ([]entities.Subscriber, error) {
	var s []entities.Subscriber
	err := db.Where("user_id = ? and id in (?)", userID, ids).Preload("Metadata").Find(&s).Error
	return s, err
}

// GetSubscriberByEmail returns the subscriber by the given email and user id
func (db *store) GetSubscriberByEmail(email string, userID int64) (*entities.Subscriber, error) {
	var s = new(entities.Subscriber)
	err := db.Where("user_id = ? and email = ?", userID, email).Preload("Metadata").Find(s).Error
	return s, err
}

// GetSubscribersByListID fetches subscribers by user id and list id, and populates the pagination obj
func (db *store) GetSubscribersByListID(listID, userID int64, p *pagination.Pagination) {
	var l = &entities.List{Id: listID}
	var subs []entities.Subscriber

	db.Model(&l).Offset(p.Offset).Limit(p.PerPage).Where("user_id = ?", 5).Association("Subscribers").Find(&subs)
	p.SetTotal(uint64(db.Model(&l).Where("user_id = ?", 5).Association("Subscribers").Count()))

	for _, t := range subs {
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

// DeleteSubscriber deletes an existing subscriber from the database along with all his metadata.
func (db *store) DeleteSubscriber(id, userID int64) error {
	s, err := db.GetSubscriber(id, userID)
	if err != nil {
		return err
	}

	var meta []entities.SubscriberMetadata

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Where("subscriber_id = ?", id).Delete(meta).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(s).Association("Lists").Clear().Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("user_id = ?", userID).Delete(s).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
