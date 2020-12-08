package s3

import (
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

// s3storage implements the S3Storage interface
type s3storage struct {
	s3client s3iface.S3API
}

// New creates a database connection and returns a new S3Storage
func New(client s3iface.S3API) S3Storage {
	return &s3storage{
		s3client: client,
	}
}

