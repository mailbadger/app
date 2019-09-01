package storage

import (
	"time"

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
	var reverse bool
	var prevID, nextID int64

	query := db.Where("user_id = ?", userID).Limit(p.PerPage).Order("created_at desc, id desc")

	if p.EndingBefore != 0 {
		sub, err := db.GetSubscriber(p.EndingBefore, userID)
		if err != nil {
			logrus.WithFields(logrus.Fields{"ending_before": p.EndingBefore, "user_id": userID}).WithError(err).
				Error("Unable to find subscriber for pagination with ending before id.")
			return
		}

		query.Where(`(created_at > ? OR (created_at = ? AND id > ?)) AND created_at < ?`,
			sub.CreatedAt.Format(time.RFC3339Nano),
			sub.CreatedAt.Format(time.RFC3339Nano),
			sub.ID,
			time.Now().Format(time.RFC3339Nano),
		).
			Order("created_at, id", true).Find(&subs)

		// populate prev and next
		if len(subs) > 0 {
			nextID = subs[0].ID
			last, err := db.getLastSubscriber(userID)
			if err != nil {
				logrus.WithFields(logrus.Fields{"user_id": userID}).WithError(err).
					Error("Unable to find the last subscriber.")
				return
			}

			if last.ID != subs[len(subs)-1].ID {
				prevID = subs[len(subs)-1].ID
			}
		}

		reverse = true
	} else if p.StartingAfter != 0 {
		sub, err := db.GetSubscriber(p.StartingAfter, userID)
		if err != nil {
			logrus.WithFields(logrus.Fields{"starting_after": p.StartingAfter, "user_id": userID}).WithError(err).
				Error("Unable to find subscriber for pagination with starting after id.")
			return
		}
		query.Where(`(created_at < ? OR (created_at = ? AND id < ?)) AND created_at < ?`,
			sub.CreatedAt.Format(time.RFC3339Nano),
			sub.CreatedAt.Format(time.RFC3339Nano),
			sub.ID,
			time.Now().Format(time.RFC3339Nano),
		).Find(&subs)

		// populate prev and next
		if len(subs) > 0 {
			prevID = subs[0].ID
			first, err := db.getFirstSubscriber(userID)
			if err != nil {
				logrus.WithFields(logrus.Fields{"user_id": userID}).WithError(err).
					Error("Unable to find the first subscriber.")
				return
			}

			if first.ID != subs[len(subs)-1].ID {
				nextID = subs[len(subs)-1].ID
			}
		}
	} else {
		query.Find(&subs)
		if len(subs) > 0 {
			nextID = subs[len(subs)-1].ID
		}
	}

	if reverse {
		for i := len(subs) - 1; i >= 0; i-- {
			p.Append(subs[i])
		}
	} else {
		for _, s := range subs {
			p.Append(s)
		}
	}

	p.PopulateLinks(prevID, nextID)
}

func (db *store) getFirstSubscriber(userID int64) (*entities.Subscriber, error) {
	var s = new(entities.Subscriber)
	err := db.Where("user_id = ?", userID).Order("created_at, id").Limit(1).Find(s).Error
	return s, err
}

func (db *store) getLastSubscriber(userID int64) (*entities.Subscriber, error) {
	var s = new(entities.Subscriber)
	err := db.Where("user_id = ?", userID).Order("created_at desc, id desc").Limit(1).Find(s).Error
	return s, err
}

func (db *store) getFirstSubscriberBySegment(segmentID, userID int64) (*entities.Subscriber, error) {
	var sub = new(entities.Subscriber)
	var seg = entities.Segment{ID: segmentID}
	err := db.Model(&seg).
		Where("user_id = ?", userID).
		Order("created_at, id").
		Limit(1).
		Association("Subscribers").
		Find(sub).Error
	return sub, err
}

func (db *store) getLastSubscriberBySegment(segmentID, userID int64) (*entities.Subscriber, error) {
	var sub = new(entities.Subscriber)
	var seg = entities.Segment{ID: segmentID}
	err := db.Model(&seg).
		Where("user_id = ?", userID).
		Order("created_at desc, id desc").
		Limit(1).
		Association("Subscribers").
		Find(sub).Error
	return sub, err
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
func (db *store) GetSubscribersBySegmentID(segmentID, userID int64, p *pagination.Cursor) {
	var subs []entities.Subscriber

	var reverse bool
	var prevID, nextID int64

	query := db.Table("subscribers").
		Joins("INNER JOIN subscribers_segments ON subscribers_segments.subscriber_id = subscribers.id").
		Where("subscribers.user_id = ? AND subscribers_segments.segment_id = ?", userID, segmentID).
		Limit(p.PerPage).
		Order("created_at desc, id desc")

	if p.EndingBefore != 0 {
		sub, err := db.GetSubscriber(p.EndingBefore, userID)
		if err != nil {
			logrus.WithFields(logrus.Fields{"ending_before": p.EndingBefore, "user_id": userID}).WithError(err).
				Error("Unable to find subscriber for pagination with ending before id.")
			return
		}

		query.Where(`(
				subscribers.created_at > ? 
				OR (subscribers.created_at = ? AND subscribers.id > ?)
			)
			AND subscribers.created_at < ?`,
			sub.CreatedAt.Format(time.RFC3339Nano),
			sub.CreatedAt.Format(time.RFC3339Nano),
			sub.ID,
			time.Now().Format(time.RFC3339Nano),
		).Order("created_at, id", true).Find(&subs)

		// populate prev and next
		if len(subs) > 0 {
			nextID = subs[0].ID
			last, err := db.getLastSubscriberBySegment(segmentID, userID)
			if err != nil {
				logrus.WithFields(logrus.Fields{"user_id": userID}).WithError(err).
					Error("Unable to find the last subscriber.")
				return
			}

			if last.ID != subs[len(subs)-1].ID {
				prevID = subs[len(subs)-1].ID
			}
		}

		reverse = true
	} else if p.StartingAfter != 0 {
		sub, err := db.GetSubscriber(p.StartingAfter, userID)
		if err != nil {
			logrus.WithFields(logrus.Fields{"starting_after": p.StartingAfter, "user_id": userID}).WithError(err).
				Error("Unable to find subscriber for pagination with starting after id.")
			return
		}
		query.Where(`(
				subscribers.created_at < ? 
				OR (subscribers.created_at = ? AND subscribers.id < ?)
			)
			AND subscribers.created_at < ?`,
			sub.CreatedAt.Format(time.RFC3339Nano),
			sub.CreatedAt.Format(time.RFC3339Nano),
			sub.ID,
			time.Now().Format(time.RFC3339Nano),
		).Find(&subs)

		// populate prev and next
		if len(subs) > 0 {
			prevID = subs[0].ID
			first, err := db.getFirstSubscriberBySegment(segmentID, userID)
			if err != nil {
				logrus.WithFields(logrus.Fields{"user_id": userID}).WithError(err).
					Error("Unable to find the first subscriber.")
				return
			}

			if first.ID != subs[len(subs)-1].ID {
				nextID = subs[len(subs)-1].ID
			}
		}
	} else {
		query.Find(&subs)
		if len(subs) > 0 {
			nextID = subs[len(subs)-1].ID
		}
	}

	if reverse {
		for i := len(subs) - 1; i >= 0; i-- {
			p.Append(subs[i])
		}
	} else {
		for _, s := range subs {
			p.Append(s)
		}
	}

	p.PopulateLinks(prevID, nextID)
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
		Select("DISTINCT(id), name, email, created_at").
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
	return db.Model(&entities.Subscriber{}).
		Where("user_id = ? AND email = ?", userID, email).
		Update("blacklisted", true).Error
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
