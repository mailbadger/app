package storage

import (
	"github.com/news-maily/app/entities"
	"github.com/news-maily/app/utils/pagination"
)

// GetSegments fetches lists by user id, and populates the pagination obj
func (db *store) GetSegments(userID int64, p *pagination.Pagination) {
	var lists []entities.Segment
	var count uint64

	db.Offset(p.Offset).Limit(p.PerPage).Where("user_id = ?", userID).Find(&lists).Count(&count)
	p.SetTotal(count)

	for _, t := range lists {
		p.Append(t)
	}
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
	err := db.Where("user_id = ? and id = ?", userID, id).Preload("Subscribers").Find(list).Error
	return list, err
}

// GetSegmentByName returns the campaign by the given name and user id
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
func (db *store) RemoveSubscribersFromSegment(l *entities.Segment) error {
	return db.Model(l).Association("Subscribers").Clear().Error
}

// AppendSubscribers appends subscribers to the existing association.
func (db *store) AppendSubscribers(l *entities.Segment) error {
	return db.Model(l).Association("Subscribers").Append(l.Subscribers).Error
}

// DetachSubscribers deletes the subscribers association by the given subscribers list.
func (db *store) DetachSubscribers(l *entities.Segment) error {
	return db.Model(l).Association("Subscribers").Delete(l.Subscribers).Error
}
