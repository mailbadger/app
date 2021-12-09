package storage

import "github.com/mailbadger/app/entities"

func (db *store) CreateSubscriberMetrics(sm []*entities.SubscribersMetrics, job *entities.Job) error {
	tx:= db.DB.Begin()
	
	// add the queries
	
	return tx.Commit().Error
}