package subscribers

import (
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

type Service interface {
	ImportSubscribersFromFile(ctx context.Context, filename string, userID int64, segments []entities.Segment) error
	RemoveSubscribersFromFile(ctx context.Context, filename string, userID int64) error
}

type service struct {
	client s3iface.S3API
	db     storage.Storage
}

var (
	ErrInvalidColumnsNum = errors.New("importer: invalid number of columns")
	ErrInvalidFormat     = errors.New("importer: csv file not formatted properly")
)

func New(client s3iface.S3API, db storage.Storage) *service {
	return &service{client, db}
}

func (s *service) ImportSubscribersFromFile(
	ctx context.Context,
	userID int64,
	segments []entities.Segment,
	r io.Reader,
) (err error) {

	reader := csv.NewReader(r)
	header, err := reader.Read()
	if err != nil {
		if err == io.EOF {
			return fmt.Errorf("importer: empty file : %w", err)
		}

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
		if err != nil {
			if err == io.EOF {
				break
			}
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

		sub := &entities.Subscriber{
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
			sub.MetaJSON = metaJSON
		}

		err = s.db.CreateSubscriber(sub)
		if err != nil {
			return fmt.Errorf("importer: create subscriber: %w", err)
		}
	}

	return
}

func (s *service) RemoveSubscribersFromFile(
	ctx context.Context,
	filename string,
	userID int64,
	r io.ReadCloser,
) (err error) {

	defer func() {
		if cerr := r.Close(); cerr != nil {
			err = cerr
		}
	}()

	reader := csv.NewReader(r)
	header, err := reader.Read()
	if err != nil {
		if err == io.EOF {
			return fmt.Errorf("bulkremover: empty file '%s': %w", filename, err)
		}

		return fmt.Errorf("bulkremover: read header: %w", err)
	}

	if len(header) < 1 {
		return ErrInvalidColumnsNum
	}

	if strings.ToLower(header[0]) != "email" {
		return ErrInvalidFormat
	}

	for {
		line, err := reader.Read()

		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("bulkremover: read line: %w", err)
		}
		if len(line) == 0 {
			continue
		}

		email := strings.TrimSpace(line[0])
		err = s.db.DeleteSubscriberByEmail(email, userID)
		if err != nil {
			return fmt.Errorf("bulkremover: delete subscriber: %w", err)
		}
	}

	return
}
