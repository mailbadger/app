package s3

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/stretchr/testify/mock"
)

// MockS3Client structure with s3API mock for testing
type MockS3Client struct {
	mock.Mock
	s3iface.S3API
}

func (m *MockS3Client) PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	args := m.Called(input)

	var obj s3.PutObjectOutput
	objBytes, _ := json.Marshal(args.Get(0))

	// lint:ignore SA4006 we can safely ignore the error check here since it is a mock
	json.Unmarshal(objBytes, &obj)

	return &obj, args.Error(1)
}

func (m *MockS3Client) GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	args := m.Called(input)

	var obj s3.GetObjectOutput
	objBytes, _ := json.Marshal(args.Get(0))

	// lint:ignore SA4006 we can safely ignore the error check here since it is a mock
	json.Unmarshal(objBytes, &obj)

	return &obj, args.Error(1)
}

func (m *MockS3Client) DeleteObject(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
	args := m.Called(input)

	var obj s3.DeleteObjectOutput
	objBytes, _ := json.Marshal(args.Get(0))

	// lint:ignore SA4006 we can safely ignore the error check here since it is a mock
	json.Unmarshal(objBytes, &obj)

	return &obj, args.Error(1)
}
