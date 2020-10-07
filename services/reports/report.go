package reports

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/services/exporters"
	"github.com/mailbadger/app/storage"
)

// ReportService represents all report functionalities
type ReportService interface {
	IsAnotherReportRunning(*gin.Context, int64) bool
	IsLimitExceeded(*gin.Context, int64) bool
	GenerateExportReport(*gin.Context, *entities.Report)
	GenerateCSV(*gin.Context)
	GenerateFilename(*gin.Context, int64, string, time.Time) string
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
	return err != nil
}

func (r reportService) IsLimitExceeded(c *gin.Context, userID int64) bool {
	panic("implement me")
}

func (r reportService) GenerateExportReport(c *gin.Context, report *entities.Report) {

}

func (r reportService) GenerateCSV(c *gin.Context) {
	panic("implement me")
}

func (r reportService) GenerateFilename(context *gin.Context, userID int64, resource string, timestamp time.Time) string {
	return fmt.Sprintf("/reports/%d/%s_%s",userID,resource,timestamp.Format("2006-01-02 15:04:05"))
}
