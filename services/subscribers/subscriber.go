package subscribers

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/jinzhu/gorm"

	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/storage"
)

type SubscriberService interface {
	ImportSubscribersFromFile(ctx context.Context, filename string, userID int64, segments []entities.Segment) error
	RemoveSubscribersFromFile(ctx context.Context, filename string, userID int64) error
	DeactivateSubscriber(ctx context.Context, userID int64, email string) error
}

type subscriberService struct {
	client s3iface.S3API
}

var (
	ErrInvalidColumnsNum = errors.New("invalid number of columns")
	ErrInvalidFormat     = errors.New("csv file not formatted properly")
)

func NewSubscriberService(client s3iface.S3API) *subscriberService {
	return &subscriberService{client}
}

func (s *subscriberService) ImportSubscribersFromFile(
	ctx context.Context,
	filename string,
	userID int64,
	segments []entities.Segment,
) (err error) {
	res, err := s.client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("FILES_BUCKET")),
		Key:    aws.String(fmt.Sprintf("subscribers/import/%d/%s", userID, filename)),
	})
	if err != nil {
		return fmt.Errorf("importer: get object: %w", err)
	}
	defer func() {
		if cerr := res.Body.Close(); cerr != nil {
			err = cerr
		}
	}()

	reader := csv.NewReader(res.Body)
	header, err := reader.Read()
	if err == io.EOF {
		return fmt.Errorf("importer: empty file '%s': %w", filename, err)
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

func (s *subscriberService) RemoveSubscribersFromFile(
	ctx context.Context,
	filename string,
	userID int64,
) (err error) {
	res, err := s.client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("FILES_BUCKET")),
		Key:    aws.String(fmt.Sprintf("subscribers/remove/%d/%s", userID, filename)),
	})
	if err != nil {
		return fmt.Errorf("bulkremover: get object: %w", err)
	}
	defer func() {
		if cerr := res.Body.Close(); cerr != nil {
			err = cerr
		}
	}()

	reader := csv.NewReader(res.Body)
	header, err := reader.Read()
	if err == io.EOF {
		return fmt.Errorf("bulkremover: empty file '%s': %w", filename, err)
	}
	if err != nil {
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
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("bulkremover: read line: %w", err)
		}
		if len(line) == 0 {
			continue
		}

		email := strings.TrimSpace(line[0])
		err = storage.DeleteSubscriberByEmail(ctx, email, userID)
		if err != nil {
			return fmt.Errorf("bulkremover: delete subscriber: %w", err)
		}
	}

	return
}

func (s *subscriberService) DeactivateSubscriber(ctx context.Context, userID int64, email string) error {

	// todo set status to inactive


	// todo insert event log for unsubscribe



	return nil

}
