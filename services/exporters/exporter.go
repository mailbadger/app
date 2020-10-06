package exporters

import (
	"github.com/gin-gonic/gin"

	"github.com/mailbadger/app/logger"
)

// Exporter represents type for creating exporters for different resource
type Exporter interface {
	Export(c *gin.Context)
}

func NewExporter(resource string) Exporter {
	switch resource {
	case "subscriptions":
		return NewSubscriptionExporter()
	default:
		return newDefaultExporter()
	}
}

type defaultExporter struct {
}

func newDefaultExporter() *defaultExporter {
	return &defaultExporter{}
}

func (de defaultExporter) Export(c *gin.Context) {
	// TODO discuss this
	logger.From(c).Errorf("Something went wrong")
}
