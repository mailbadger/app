package exporters

import (
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/gin-gonic/gin"

	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/logger"
)

// Exporter represents type for creating exporters for different resource
type Exporter interface {
	Export(*gin.Context, *entities.Report) error
}

func NewExporter(resource string, s3 s3iface.S3API) Exporter {
	switch resource {
	case "subscriptions":
		return NewSubscriptionExporter(s3)
	default:
		return newDefaultExporter()
	}
}

type defaultExporter struct {
}

func newDefaultExporter() *defaultExporter {
	return &defaultExporter{}
}

func (de defaultExporter) Export(c *gin.Context, report *entities.Report) error {
	// TODO discuss this
	logger.From(c).Errorf("Something went wrong")
	return nil
}
