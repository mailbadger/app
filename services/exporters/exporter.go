package exporters

import (
	"context"

	"github.com/mailbadger/app/entities"
)

// Exporter represents type for creating exporters for different resource
type Exporter interface {
	Export(c context.Context, userID int64, report *entities.Report, bucket string) error
}
