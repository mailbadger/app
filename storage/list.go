package storage

import (
	"github.com/FilipNikolovski/news-maily/entities"
	"github.com/FilipNikolovski/news-maily/utils/pagination"
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

// GetList returns the list by the given id and user id
func (db *store) GetList(id, userID int64) (*entities.List, error) {
	var list = new(entities.List)
	err := db.Where("user_id = ? and id = ?", userID, id).Find(list).Error
	return list, err
}

// CreateList creates a new list in the database.
func (db *store) CreateList(l *entities.List) error {
	return db.Create(l).Error
}

// UpdateList edits an existing list in the database.
func (db *store) UpdateList(l *entities.List) error {
	return db.Where("id = ? and user_id = ?", l.Id, l.UserId).Save(l).Error
}

// DeleteList deletes an existing list from the database.
func (db *store) DeleteList(id, userID int64) error {
	return db.Where("user_id = ?", userID).Delete(entities.List{Id: id}).Error
}
