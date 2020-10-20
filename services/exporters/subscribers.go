package exporters

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"

	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/storage"
)

type SubscribersExporter struct {
	S3 s3iface.S3API
}

func NewSubscriptionExporter(s3 s3iface.S3API) *SubscribersExporter {
	return &SubscribersExporter{
		S3: s3,
	}
}

func (se *SubscribersExporter) Export(c context.Context, report *entities.Report) error {
	var (
		nextID      int64
		limit       int64 = 1000
	)

	for {
		subscribers, err := storage.GetSubscribersByUserID(c, report.UserID, nextID, limit)
		if err != nil {
			return fmt.Errorf("get subscribers: %w", err)
		}

		// TODO write these 1000 subscribers into the csv


		if len(subscribers) < 1000 {
			break
		}

		nextID = subscribers[len(subscribers)-1].ID
	}

	// i will comment this code for now
	/*_, err = se.S3.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("AWS_S3_BUCKET")),
		Key:    aws.String(fmt.Sprintf("subscribers/export/%d/%s", report.UserID, report.FileName)),
		Body:   bytes.NewReader(reportDataBytes),
	})
	if err != nil {
		return fmt.Errorf("put object: %w", err)
	}*/

	return nil
}
