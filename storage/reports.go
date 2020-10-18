package storage

import (
	"time"

	"github.com/mailbadger/app/entities"
)

// CreateReport creates a report.
func (db *store) CreateReport(r *entities.Report) error {
	return db.Create(r).Error
}

// UpdateReport edits an existing report in the database.
func (db *store) UpdateReport(r *entities.Report) error {
	return db.Where("id = ? and user_id = ?", r.ID, r.UserID).Save(r).Error
}

// GetReportByFilename returns the report by the given file name and user id
func (db *store) GetReportByFilename(filename string, userID int64) (*entities.Report, error) {
	var report = new(entities.Report)
	err := db.Where("user_id = ? and file_name = ?", userID, filename).Find(report).Error
	return report, err
}

func (db *store) GetRunningReportForUser(userID int64) (*entities.Report, error) {
	var report = new(entities.Report)
	err := db.Where("user_id = ? and status = ?", userID, entities.StatusInProgress).Find(report).Error
	return report, err
}

// GetNumberOfReportsForDateTime returns number of reports for user id and datetime.
func (db *store) GetNumberOfReportsForDate(userID int64, time time.Time) (int64, error) {
	var count int64
	err := db.Model(entities.Report{}).Where("user_id = ? and DATE(created_at) = DATE(?)", userID, time).Count(&count).Error
	return count, err
}
