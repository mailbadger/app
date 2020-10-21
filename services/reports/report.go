package reports

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/services/exporters"
	"github.com/mailbadger/app/storage"
)

const (
	reportTypeExport = "export"
)

var (
	ErrAnotherReportRunning = errors.New("another report running")
	ErrLimitReached         = errors.New("you reached the limit")
)

// ReportService represents all report functionalities
type ReportService interface {
	GenerateExportReport(context.Context, *entities.Report) error
	CreateExportReport(*gin.Context, int64, string, string, time.Time) (*entities.Report, error)
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

func (r *reportService) GenerateExportReport(c context.Context, report *entities.Report) error {
	err := r.exporter.Export(c, report)
	if err != nil {
		return fmt.Errorf("export: %w", err)
	}

	return nil
}

func (r *reportService) CreateExportReport(c *gin.Context, userID int64, resource, note string, date time.Time) (*entities.Report, error) {
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
		FileName: generateFilename(c, userID, resource, date),
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

func generateFilename(context *gin.Context, userID int64, resource string, timestamp time.Time) string {
	return fmt.Sprintf("/reports/%d/%s_%s", userID, resource, timestamp.Format("2006-01-02 15:04:05"))
}

// isAnotherReportRunning returns true if there is report in progress for a user or false if all are done
func isAnotherReportRunning(c *gin.Context, userID int64) bool {
	_, err := storage.GetRunningReportForUser(c, userID)
	return err != nil
}

func isLimitExceeded(c *gin.Context, userID int64, time time.Time) (bool, error) {
	n, err := storage.GetNumberOfReportsForDate(c, userID, time)
	if err != nil {
		return false, err
	}

	if n > 100 {
		return true, nil
	}
	return false, nil
}
