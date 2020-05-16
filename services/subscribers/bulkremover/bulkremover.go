package bulkremover

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/mailbadger/app/storage"
)

type SubscribersBulkRemover interface {
	RemoveSubscribersFromFile(ctx context.Context, filename string, userID int64) error
}

type s3Remover struct {
	client s3iface.S3API
}

var (
	ErrInvalidColumnsNum = errors.New("bulkremover: invalid number of columns")
	ErrInvalidFormat     = errors.New("bulkremover: csv file not formatted properly")
)

func NewS3SubscribersBulkRemover(client s3iface.S3API) *s3Remover {
	return &s3Remover{client}
}

func (svc *s3Remover) RemoveSubscribersFromFile(
	ctx context.Context,
	filename string,
	userID int64,
) error {
	res, err := svc.client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("AWS_S3_BUCKET")),
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
	return nil
}
