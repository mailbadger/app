package reports

import (
	"github.com/gin-gonic/gin"

	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/services/exporters"
	"github.com/mailbadger/app/storage"
)

type ReportService interface {
	IsAnotherReportRunning(*gin.Context, int64) bool
	IsLimitExceeded() bool
	GenerateReport() (*entities.Report, error)
	UpdateReport() (*entities.Report, error)
	GenerateCSV()
}

type reportService struct {
	exporter exporters.Exporter
}

func NewReportService(exporter exporters.Exporter) ReportService {
	return &reportService{
		exporter: exporter,
	}
}

func (r reportService) IsAnotherReportRunning(c *gin.Context, userID int64) bool {
	_, err := storage.GetRunningReportForUser(c, userID)
	if err != nil {
		return true
	}
	return false
}

func (r reportService) IsLimitExceeded() bool {
	panic("implement me")
}

func (r reportService) GenerateReport() (*entities.Report, error) {
	panic("implement me")
}

func (r reportService) UpdateReport() (*entities.Report, error) {
	panic("implement me")
}

func (r reportService) GenerateCSV() {
	panic("implement me")
}
