package storage

import "github.com/mailbadger/app/entities"

// GetJobByName returns the job by the given name
func (db store) GetJobByName(name string) (*entities.Job, error) {
	var job = new(entities.Job)
	err := db.Where("name = ?", name).Find(job).Error
	return job, err
}

// UpdateJob edits an existing job in the database.
func (db *store) UpdateJob(job *entities.Job) error {
	return db.Where("id = ? and name = ?", job.ID, job.Name).Save(job).Error
}
