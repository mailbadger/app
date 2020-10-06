package reports

import (
	"github.com/gin-gonic/gin"

	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/services/exporters"
	"github.com/mailbadger/app/storage"
)

// ReportService represents all report functionalities
type ReportService interface {
	IsAnotherReportRunning(*gin.Context, int64) bool
	IsLimitExceeded(*gin.Context, int64) bool
	GenerateReport(*gin.Context) (*entities.Report, error)
	UpdateReport(*gin.Context) (*entities.Report, error)
	GenerateCSV(*gin.Context)
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

// IsAnotherReportRunning returns true if there is report in progress for a user or false if all are done
func (r reportService) IsAnotherReportRunning(c *gin.Context, userID int64) bool {
	_, err := storage.GetRunningReportForUser(c, userID)
	if err != nil {
		return true
	}
	return false
}

func (r reportService) IsLimitExceeded(c *gin.Context, userID int64) bool {
	panic("implement me")
}

func (r reportService) GenerateReport(c *gin.Context) (*entities.Report, error) {
	panic("implement me")
}

func (r reportService) UpdateReport(c *gin.Context) (*entities.Report, error) {
	panic("implement me")
}

func (r reportService) GenerateCSV(c *gin.Context) {
	panic("implement me")
}
