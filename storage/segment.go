package storage

import (
	"github.com/mailbadger/app/entities"
)

// GetSegments fetches lists by user id, and populates the pagination obj
func (db *store) GetSegments(userID int64, p *PaginationCursor) error {
	p.SetCollection(&[]entities.SegmentWithTotalSubs{})
	p.SetResource("segments")

	query := db.Table(p.Resource).
		Select("segments.*, (?) as subscribers_in_segment",
			db.Select("count(*)").
				Table("subscribers_segments").
				Where("segment_id = segments.id"),
		).
		Where("user_id = ?", userID).
		Order("created_at desc, id desc").
		Limit(int(p.PerPage))

	p.SetQuery(query)

	return db.Paginate(p, userID)
}

// GetTotalSegments fetches the total count by user id
func (db *store) GetTotalSegments(userID int64) (int64, error) {
	var count int64
	err := db.Model(entities.Segment{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// GetSegmentsByIDs fetches lists by user id and the given ids
func (db *store) GetSegmentsByIDs(userID int64, ids []int64) ([]entities.Segment, error) {
	var lists []entities.Segment

	err := db.Where("user_id = ? AND id IN (?)", userID, ids).Find(&lists).Error

	return lists, err
}

// GetSegment returns the list by the given id and user id
func (db *store) GetSegment(id, userID int64) (*entities.Segment, error) {
	var seg = new(entities.Segment)
	err := db.Where("user_id = ? and id = ?", userID, id).Find(seg).Error
	return seg, err
}

// GetSegmentByName returns the segment by the given name and user id
func (db *store) GetSegmentByName(name string, userID int64) (*entities.Segment, error) {
	var seg = new(entities.Segment)
	err := db.Where("user_id = ? and name = ?", userID, name).Find(seg).Error
	return seg, err
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
	l := &entities.Segment{Model: entities.Model{ID: id}, UserID: userID}
	if err := db.RemoveSubscribersFromSegment(l); err != nil {
		return err
	}

	return db.Delete(&l).Error
}

// RemoveSubscribersFromSegment clears the subscribers association.
func (db *store) RemoveSubscribersFromSegment(s *entities.Segment) error {
	return db.Model(s).Association("Subscribers").Clear()
}

// AppendSubscribers appends segscribers to the existing association.
func (db *store) AppendSubscribers(s *entities.Segment) error {
	return db.Model(s).Association("Subscribers").Append(s.Subscribers)
}

// DetachSubscribers deletes the subscribers association by the given subscribers list.
func (db *store) DetachSubscribers(s *entities.Segment) error {
	return db.Model(s).Association("Subscribers").Delete(s.Subscribers)
}

// DeleteAllSegmentsForUser deletes all subscribers for user
func (db *store) DeleteAllSegmentsForUser(userID int64) error {
	return db.Where("user_id = ?", userID).Delete(&entities.Segment{}).Error
}
