package exporters

import "github.com/gin-gonic/gin"

type Exporter interface {
	Export(c *gin.Context)
}

func NewExporter(resource string) Exporter {
	switch resource {
	case "subscriptions":
		return NewSubscriptionExporter()
	}
}