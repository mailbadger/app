package storage

import "github.com/mailbadger/app/entities"

func (db *store) GetBoundariesByType(t string) (*entities.Boundaries, error) {
	var b = new(entities.Boundaries)
	err := db.Where("type = ?", t).Find(b).Error
	return b, err
}
