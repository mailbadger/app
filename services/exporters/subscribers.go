package exporters

import "github.com/gin-gonic/gin"

type SubscribersExporter struct {
}

func NewSubscriptionExporter() *SubscribersExporter {
	return &SubscribersExporter{}
}

func (se *SubscribersExporter) Export (c *gin.Context) {
}