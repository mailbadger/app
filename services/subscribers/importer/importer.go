package importer

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/jinzhu/gorm"

	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/storage"
)

const ActionImport = "import"

type SubscribersImporter interface {
	ImportSubscribersFromFile(ctx context.Context, filename string, userID int64, segments []entities.Segment) error
}

type s3Importer struct {
	client s3iface.S3API
}

var (
	ErrInvalidColumnsNum = errors.New("importer: invalid number of columns")
	ErrInvalidFormat     = errors.New("importer: csv file not formatted properly")
)

func NewS3SubscribersImporter(client s3iface.S3API) *s3Importer {
	return &s3Importer{client}
}

func (i *s3Importer) ImportSubscribersFromFile(
	ctx context.Context,
	userID int64,
	segments []entities.Segment,
	body []byte,
) (err error) {

	r := bytes.NewReader(body)

	reader := csv.NewReader(r)
	header, err := reader.Read()
	if err == io.EOF {
		return fmt.Errorf("importer: empty file: %w", err)
	}
	if err != nil {
		return fmt.Errorf("importer: read header: %w", err)
	}

	if len(header) < 2 {
		return ErrInvalidColumnsNum
	}

	if strings.ToLower(header[0]) != "email" || strings.ToLower(header[1]) != "name" {
		return ErrInvalidFormat
	}

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("importer: read line: %w", err)
		}
		if len(line) < 2 {
			continue
		}
		email := strings.TrimSpace(line[0])
		name := strings.TrimSpace(line[1])

		_, err = storage.GetSubscriberByEmail(ctx, email, userID)
		if err == nil {
			continue
		} else if !gorm.IsRecordNotFoundError(err) {
			return fmt.Errorf("importer: get subscriber by email: %w", err)
		}

		s := &entities.Subscriber{
			UserID:   userID,
			Email:    email,
			Name:     name,
			Segments: segments,
			Active:   true,
		}

		if len(line) > 2 {
			meta := make(map[string]string, len(line)-2)
			keys := header[2:]
			for i, m := range line[2:] {
				meta[keys[i]] = m
			}
			metaJSON, err := json.Marshal(meta)
			if err != nil {
				return fmt.Errorf("importer: marshal metadata: %w", err)
			}
			s.MetaJSON = metaJSON
		}

		err = storage.CreateSubscriber(ctx, s)
		if err != nil {
			return fmt.Errorf("importer: create subscriber: %w", err)
		}
	}

	return
}
