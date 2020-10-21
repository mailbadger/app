package reports

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/services/exporters"
	"github.com/mailbadger/app/storage"
)

const (
	reportTypeExport = "export"
)

var (
	ErrAnotherReportRunning = errors.New("another report running")
	ErrLimitReached         = errors.New("report limit reached")
)

// ReportService represents all report functionalities
type ReportService interface {
	GenerateExportReport(c context.Context, report *entities.Report) error
	CreateExportReport(c context.Context, userID int64, resource string, note string, date time.Time) (*entities.Report, error)
}

type reportService struct {
	exporter exporters.Exporter
}

// NewReportService represents constructor for ReportService
func NewReportService(exporter exporters.Exporter) ReportService {
	return &reportService{
		exporter: exporter,
	}
}

// GenerateExportReport starts the resources export method
func (r *reportService) GenerateExportReport(c context.Context, report *entities.Report) error {
	err := r.exporter.Export(c, report)
	if err != nil {
		return fmt.Errorf("export: %w", err)
	}

	return nil
}

// CreateExportReport creates export report
func (r *reportService) CreateExportReport(c context.Context, userID int64, resource, note string, date time.Time) (*entities.Report, error) {
	if isAnotherReportRunning(c, userID) {
		return nil, ErrAnotherReportRunning
	}

	limit, err := isLimitExceeded(c, userID, date)
	if err != nil {
		return nil, fmt.Errorf("is limit exceeded check error: %w", err)
	}
	if limit {
		return nil, ErrLimitReached
	}

	report := &entities.Report{
		UserID:   userID,
		Resource: resource,
		FileName: generateFilename(userID, resource, date),
		Type:     reportTypeExport,
		Status:   entities.StatusInProgress,
		Note:     note,
	}

	err = storage.CreateReport(c, report)
	if err != nil {
		return nil, fmt.Errorf("creaate report: %w", err)
	}

	return report, nil
}

func generateFilename(userID int64, resource string, timestamp time.Time) string {
	return fmt.Sprintf("/reports/%d/%s_%s", userID, resource, timestamp.Format("2006-01-02 15:04:05"))
}

// isAnotherReportRunning returns true if there is report in progress for a user or false if all are done
func isAnotherReportRunning(c context.Context, userID int64) bool {
	_, err := storage.GetRunningReportForUser(c, userID)
	return err != nil
}

// isLimitExceeded returns true if there are less than 100 reports per day
func isLimitExceeded(c context.Context, userID int64, time time.Time) (bool, error) {
	n, err := storage.GetNumberOfReportsForDate(c, userID, time)
	if err != nil {
		return false, err
	}

	return n > 100, nil
}
