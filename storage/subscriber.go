package storage

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"

	"github.com/news-maily/app/entities"
	"github.com/news-maily/app/utils/pagination"
)

// GetSubscribers fetches subscribers by user id, and populates the pagination obj
func (db *store) GetSubscribers(
	userID int64,
	p *pagination.Cursor,
) {
	var subs []entities.Subscriber
	var query *gorm.DB
	var reverse bool

	if p.EndingBefore != 0 {
		sub, err := db.GetSubscriber(p.EndingBefore, userID)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"ending_before": p.EndingBefore,
				"user_id":       userID,
			}).WithError(err).Error("Unable to find subscriber for pagination with ending before id.")
			return
		}

		query = db.Where(`user_id = ?
			AND (created_at > ? OR (created_at = ? AND id > ?))
			AND created_at < ?`, userID, sub.CreatedAt, sub.CreatedAt, sub.ID, time.Now()).
			Order("created_at, id")

		reverse = true
	} else if p.StartingAfter != 0 {
		sub, err := db.GetSubscriber(p.StartingAfter, userID)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"starting_after": p.StartingAfter,
				"user_id":        userID,
			}).WithError(err).Error("Unable to find subscriber for pagination with starting after id.")
			return
		}
		query = db.Where(`user_id = ?
			AND (created_at < ? OR (created_at = ? AND id < ?))
			AND created_at < ?`, userID, sub.CreatedAt, sub.CreatedAt, sub.ID, time.Now()).
			Order("created_at desc, id desc")
	} else {
		query = db.Where("user_id = ?", userID).Order("created_at desc, id desc")
	}

	query.Limit(p.PerPage).Find(&subs)

	if reverse {
		for i := len(subs) - 1; i >= 0; i-- {
			p.Append(subs[i])
		}
	} else {
		for _, s := range subs {
			p.Append(s)
		}
	}

	p.PopulateLinks(1, 5)
}

// GetSubscriber returns the subscriber by the given id and user id
func (db *store) GetSubscriber(id, userID int64) (*entities.Subscriber, error) {
	var s = new(entities.Subscriber)
	err := db.Where("user_id = ? and id = ?", userID, id).Find(s).Error
	return s, err
}

// GetSubscribersByIDs returns the subscriber by the given id and user id
func (db *store) GetSubscribersByIDs(ids []int64, userID int64) ([]entities.Subscriber, error) {
	var s []entities.Subscriber
	err := db.Where("user_id = ? and id in (?)", userID, ids).Find(&s).Error
	return s, err
}

// GetSubscriberByEmail returns the subscriber by the given email and user id
func (db *store) GetSubscriberByEmail(email string, userID int64) (*entities.Subscriber, error) {
	var s = new(entities.Subscriber)
	err := db.Where("user_id = ? and email = ?", userID, email).Find(s).Error
	return s, err
}

// GetSubscribersBySegmentID fetches subscribers by user id and list id, and populates the pagination obj
func (db *store) GetSubscribersBySegmentID(listID, userID int64, p *pagination.Cursor) {
	var l = &entities.Segment{ID: listID}
	var subs []entities.Subscriber

	db.Model(&l).Where("user_id = ?", userID).Association("Subscribers").Find(&subs)

	for _, t := range subs {
		p.Append(t)
	}
}

// GetAllSubscribersBySegmentID fetches all subscribers by user id and list id
func (db *store) GetAllSubscribersBySegmentID(listID, userID int64) ([]entities.Subscriber, error) {
	var l = &entities.Segment{ID: listID}
	var subs []entities.Subscriber
	err := db.Model(&l).Where("user_id = ?", userID).Association("Subscribers").Find(&subs).Error
	return subs, err
}

// GetDistinctSubscribersBySegmentIDs fetches all distinct subscribers by user id and list ids
func (db *store) GetDistinctSubscribersBySegmentIDs(
	listIDs []int64,
	userID int64,
	blacklisted, active bool,
	timestamp time.Time,
	nextID int64,
	limit int64,
) ([]entities.Subscriber, error) {
	if limit == 0 {
		limit = 1000
	}

	var subs []entities.Subscriber

	err := db.Table("subscribers").
		Select("DISTINCT(id), name, email").
		Joins("INNER JOIN subscribers_segments ON subscribers_segments.subscriber_id = subscribers.id").
		Where(`
			subscribers_segments.segment_id IN (?)
			AND subscribers.user_id = ? 
			AND subscribers.blacklisted = ? 
			AND subscribers.active = ?
			AND (created_at > ? OR (created_at = ? AND id > ?))
			AND created_at < ?`, listIDs, userID, blacklisted, active, timestamp, timestamp, nextID, time.Now()).
		Order("created_at, id").
		Limit(limit).
		Find(&subs).Error

	return subs, err
}

// CreateSubscriber creates a new subscriber in the database.
func (db *store) CreateSubscriber(s *entities.Subscriber) error {
	return db.Create(s).Error
}

// UpdateSubscriber edits an existing subscriber in the database.
func (db *store) UpdateSubscriber(s *entities.Subscriber) error {
	return db.Where("id = ? and user_id = ?", s.ID, s.UserID).Save(s).Error
}

func (db *store) BlacklistSubscriber(userID int64, email string) error {
	return db.Model(&entities.Subscriber{}).Where("user_id = ? AND email = ?", userID, email).Update("blacklisted", true).Error
}

// DeleteSubscriber deletes an existing subscriber from the database along with all his metadata.
func (db *store) DeleteSubscriber(id, userID int64) error {
	s, err := db.GetSubscriber(id, userID)
	if err != nil {
		return err
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Model(s).Association("Segments").Clear().Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("user_id = ?", userID).Delete(s).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
