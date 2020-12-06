package exporters

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"

	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/storage"
)

type SubscribersExporter struct {
	S3 s3iface.S3API
}

func NewSubscribersExporter(s3 s3iface.S3API) *SubscribersExporter {
	return &SubscribersExporter{
		S3: s3,
	}
}

func (se *SubscribersExporter) Export(c context.Context, userID int64, report *entities.Report) error {
	var (
		err    error
		nextID int64
		limit  int64 = 1000

		buf bytes.Buffer
	)

	writer := csv.NewWriter(&buf)

	// writing headers
	// change this function to change the headers
	err = writeHeaders(writer)
	if err != nil {
		return fmt.Errorf("write headers: %w", err)
	}

	for {
		subscribers, err := storage.SeekSubscribersByUserID(c, userID, nextID, limit)
		if err != nil {
			return fmt.Errorf("get subscribers: %w", err)
		}

		// writing subscribers
		err = writeSubscribers(writer, subscribers)
		if err != nil {
			return fmt.Errorf("write %d subscribers with id greater than %d: %w", limit, nextID, err)
		}

		if len(subscribers) < int(limit) {
			break
		}

		nextID = subscribers[len(subscribers)-1].ID
	}

	writer.Flush()

	_, err = se.S3.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("AWS_S3_BUCKET")),
		Key:    aws.String(fmt.Sprintf("subscribers/export/%d/%s", userID, report.FileName)),
		Body:   bytes.NewReader(buf.Bytes()),
	})
	if err != nil {
		return fmt.Errorf("put object: %w", err)
	}

	return nil
}

// writeHeaders writes headers to csv
func writeHeaders(writer *csv.Writer) error {
	return writer.Write([]string{
		"Name",
		"Email",
		"User ID",
		"Segments",
		"Active",
		"Metadata",
		"Blacklisted",
		"Created At",
	})
}

// writeSubscribers writes the given subscribers into the csv
func writeSubscribers(writer *csv.Writer, subscribers []entities.Subscriber) error {
	for _, s := range subscribers {
		_, err := s.GetMetadata()
		if err != nil {
			return fmt.Errorf("get metadata: %w", err)
		}

		formattedSegments, err := formatSegments(s.Segments)
		if err != nil {
			return fmt.Errorf("format segments: %w", err)
		}

		formatMetadata, err := formatMetadata(s.Metadata)
		if err != nil {
			return fmt.Errorf("format metadata: %w", err)
		}

		err = writer.Write([]string{
			s.Name,
			s.Email,
			strconv.FormatInt(s.UserID, 10),
			formattedSegments,
			strconv.FormatBool(s.Active),
			formatMetadata,
			strconv.FormatBool(s.Blacklisted),
			s.GetCreatedAt().Format("2006-01-02 15:04:05"),
		})
		if err != nil {
			return fmt.Errorf("write: %w", err)
		}
	}

	return nil
}

// formatSegments returns segments divided by ;
func formatSegments(segments []entities.Segment) (string, error) {
	var b strings.Builder

	if len(segments) == 0 {
		return "", nil
	}

	for _, s := range segments {
		_, err := fmt.Fprintf(&b, "%s; ", s.Name)
		if err != nil {
			return "", fmt.Errorf("fprintf: %w", err)
		}
	}

	s := b.String()
	s = s[:b.Len()-2] // remove trailing "; "
	return s, nil
}

// formatMetadata returns metadata formatted in key = value pairs divided by ;
func formatMetadata(metadata map[string]string) (string, error) {
	var b strings.Builder

	if len(metadata) == 0 {
		return "", nil
	}

	for k, v := range metadata {
		_, err := fmt.Fprintf(&b, "%s = %s; ", k, v)
		if err != nil {
			return "", fmt.Errorf("fprintf: %w", err)
		}
	}

	s := b.String()
	s = s[:b.Len()-2] // remove trailing "; "
	return s, nil
}
