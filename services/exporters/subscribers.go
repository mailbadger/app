package exporters

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/gocarina/gocsv"

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
	p := storage.NewPaginationCursor("", 10)
	err := storage.GetSubscribers(c, report.UserID, p)
	if err != nil {
		return fmt.Errorf("get subscribers: %w", err)
	}

	subscribers := p.Collection.([]*entities.Subscriber)

	//Setting the CSV writer settings
	gocsv.SetCSVWriter(func(out io.Writer) *gocsv.SafeCSVWriter {
		writer := csv.NewWriter(out)
		writer.Comma = ','
		return gocsv.NewSafeCSVWriter(writer)
	})

	reportDataBytes, err := gocsv.MarshalBytes(subscribers)
	if err != nil {
		return fmt.Errorf("marshal csv bytes: %w", err)
	}

	_, err = se.S3.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("AWS_S3_BUCKET")),
		Key:    aws.String(fmt.Sprintf("subscribers/export/%d/%s", report.UserID, report.FileName)),
		Body:   bytes.NewReader(reportDataBytes),
	})
	if err != nil {
		return fmt.Errorf("put object: %w", err)
	}

	return nil
}
