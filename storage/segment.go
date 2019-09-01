package storage

import (
	"time"

	"github.com/news-maily/app/entities"
	"github.com/news-maily/app/utils/pagination"
	"github.com/sirupsen/logrus"
)

// GetSegments fetches lists by user id, and populates the pagination obj
func (db *store) GetSegments(userID int64, p *pagination.Cursor) {
	var seg []entities.Segment

	var reverse bool
	var prevID, nextID int64

	query := db.Where("user_id = ?", userID).Limit(p.PerPage).Order("created_at desc, id desc")

	if p.EndingBefore != 0 {
		s, err := db.GetSegment(p.EndingBefore, userID)
		if err != nil {
			logrus.WithFields(logrus.Fields{"ending_before": p.EndingBefore, "user_id": userID}).WithError(err).
				Error("Unable to find segment for pagination with ending before id.")
			return
		}

		query.Where(`(created_at > ? OR (created_at = ? AND id > ?)) AND created_at < ?`,
			s.CreatedAt.Format(time.RFC3339Nano),
			s.CreatedAt.Format(time.RFC3339Nano),
			s.ID,
			time.Now().Format(time.RFC3339Nano),
		).
			Order("created_at, id", true).Find(&seg)

		// populate prev and next
		if len(seg) > 0 {
			nextID = seg[0].ID
			last, err := db.getLastSegment(userID)
			if err != nil {
				logrus.WithFields(logrus.Fields{"user_id": userID}).WithError(err).
					Error("Unable to find the last segment.")
				return
			}

			if last.ID != seg[len(seg)-1].ID {
				prevID = seg[len(seg)-1].ID
			}
		}

		reverse = true
	} else if p.StartingAfter != 0 {
		sub, err := db.GetSegment(p.StartingAfter, userID)
		if err != nil {
			logrus.WithFields(logrus.Fields{"starting_after": p.StartingAfter, "user_id": userID}).WithError(err).
				Error("Unable to find segment for pagination with starting after id.")
			return
		}
		query.Where(`(created_at < ? OR (created_at = ? AND id < ?)) AND created_at < ?`,
			sub.CreatedAt.Format(time.RFC3339Nano),
			sub.CreatedAt.Format(time.RFC3339Nano),
			sub.ID,
			time.Now().Format(time.RFC3339Nano),
		).Find(&seg)

		// populate prev and next
		if len(seg) > 0 {
			prevID = seg[0].ID
			first, err := db.getFirstSegment(userID)
			if err != nil {
				logrus.WithFields(logrus.Fields{"user_id": userID}).WithError(err).
					Error("Unable to find the first segment.")
				return
			}

			if first.ID != seg[len(seg)-1].ID {
				nextID = seg[len(seg)-1].ID
			}
		}
	} else {
		query.Find(&seg)
		if len(seg) > 0 {
			nextID = seg[len(seg)-1].ID
		}
	}

	if reverse {
		for i := len(seg) - 1; i >= 0; i-- {
			p.Append(seg[i])
		}
	} else {
		for _, s := range seg {
			p.Append(s)
		}
	}

	p.PopulateLinks(prevID, nextID)
}

func (db *store) getFirstSegment(userID int64) (*entities.Segment, error) {
	var s = new(entities.Segment)
	err := db.Where("user_id = ?", userID).Order("created_at, id").Limit(1).Find(s).Error
	return s, err
}

func (db *store) getLastSegment(userID int64) (*entities.Segment, error) {
	var s = new(entities.Segment)
	err := db.Where("user_id = ?", userID).Order("created_at desc, id desc").Limit(1).Find(s).Error
	return s, err
}

// GetSegmentsByIDs fetches lists by user id and the given ids
func (db *store) GetSegmentsByIDs(userID int64, ids []int64) ([]entities.Segment, error) {
	var lists []entities.Segment

	err := db.Where("user_id = ? AND id IN (?)", userID, ids).Find(&lists).Error

	return lists, err
}

// GetSegment returns the list by the given id and user id
func (db *store) GetSegment(id, userID int64) (*entities.Segment, error) {
	var list = new(entities.Segment)
	err := db.Where("user_id = ? and id = ?", userID, id).Find(list).Error
	return list, err
}

// GetSegmentByName returns the segment by the given name and user id
func (db *store) GetSegmentByName(name string, userID int64) (*entities.Segment, error) {
	var list = new(entities.Segment)
	err := db.Where("user_id = ? and name = ?", userID, name).Find(list).Error
	return list, err
}

// CreateSegment creates a new list in the database.
func (db *store) CreateSegment(l *entities.Segment) error {
	return db.Create(l).Error
}

// UpdateSegment edits an existing list in the database.
func (db *store) UpdateSegment(l *entities.Segment) error {
	return db.Where("id = ? and user_id = ?", l.ID, l.UserID).Save(l).Error
}

// DeleteSegment deletes an existing list from the database and also clears the subscribers association.
func (db *store) DeleteSegment(id, userID int64) error {
	l := &entities.Segment{ID: id, UserID: userID}
	if err := db.RemoveSubscribersFromSegment(l); err != nil {
		return err
	}

	return db.Delete(&l).Error
}

// RemoveSubscribersFromSegment clears the subscribers association.
func (db *store) RemoveSubscribersFromSegment(s *entities.Segment) error {
	return db.Model(s).Association("Subscribers").Clear().Error
}

// AppendSubscribers appends subscribers to the existing association.
func (db *store) AppendSubscribers(s *entities.Segment) error {
	return db.Model(s).Association("Subscribers").Append(s.Subscribers).Error
}

// DetachSubscribers deletes the subscribers association by the given subscribers list.
func (db *store) DetachSubscribers(s *entities.Segment) error {
	return db.Model(s).Association("Subscribers").Delete(s.Subscribers).Error
}
