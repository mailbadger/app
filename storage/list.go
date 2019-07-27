package storage

import (
	"github.com/news-maily/app/entities"
	"github.com/news-maily/app/utils/pagination"
)

// GetLists fetches lists by user id, and populates the pagination obj
func (db *store) GetLists(userID int64, p *pagination.Pagination) {
	var lists []entities.List
	var count uint64

	db.Offset(p.Offset).Limit(p.PerPage).Where("user_id = ?", userID).Find(&lists).Count(&count)
	p.SetTotal(count)

	for _, t := range lists {
		p.Append(t)
	}
}

// GetListsByIDs fetches lists by user id and the given ids
func (db *store) GetListsByIDs(userID int64, ids []int64) ([]entities.List, error) {
	var lists []entities.List

	err := db.Where("user_id = ? AND id IN (?)", userID, ids).Find(&lists).Error

	return lists, err
}

// GetList returns the list by the given id and user id
func (db *store) GetList(id, userID int64) (*entities.List, error) {
	var list = new(entities.List)
	err := db.Where("user_id = ? and id = ?", userID, id).Preload("Subscribers").Find(list).Error
	return list, err
}

// GetListByName returns the campaign by the given name and user id
func (db *store) GetListByName(name string, userID int64) (*entities.List, error) {
	var list = new(entities.List)
	err := db.Where("user_id = ? and name = ?", userID, name).Find(list).Error
	return list, err
}

// CreateList creates a new list in the database.
func (db *store) CreateList(l *entities.List) error {
	return db.Create(l).Error
}

// UpdateList edits an existing list in the database.
func (db *store) UpdateList(l *entities.List) error {
	return db.Where("id = ? and user_id = ?", l.ID, l.UserID).Save(l).Error
}

// DeleteList deletes an existing list from the database and also clears the subscribers association.
func (db *store) DeleteList(id, userID int64) error {
	l := &entities.List{ID: id, UserID: userID}
	if err := db.RemoveSubscribersFromList(l); err != nil {
		return err
	}

	return db.Delete(&l).Error
}

// RemoveSubscribersFromList clears the subscribers association.
func (db *store) RemoveSubscribersFromList(l *entities.List) error {
	return db.Model(l).Association("Subscribers").Clear().Error
}

// AppendSubscribers appends subscribers to the existing association.
func (db *store) AppendSubscribers(l *entities.List) error {
	return db.Model(l).Association("Subscribers").Append(l.Subscribers).Error
}

// DetachSubscribers deletes the subscribers association by the given subscribers list.
func (db *store) DetachSubscribers(l *entities.List) error {
	return db.Model(l).Association("Subscribers").Delete(l.Subscribers).Error
}
