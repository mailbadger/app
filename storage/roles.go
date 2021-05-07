package storage

import "github.com/mailbadger/app/entities"

// GetRole fetches a role by the given name.
func (db *store) GetRole(name string) (*entities.Role, error) {
	var r = new(entities.Role)
	err := db.Where("name = ?", name).Find(r).Error
	return r, err
}
