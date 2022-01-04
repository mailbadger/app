package sqs

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockPublisher structure with sqs api mock for testing
type MockPublisher struct {
	mock.Mock
}

func (m *MockPublisher) SendMessage(ctx context.Context, queueURL *string, body []byte) error {
	args := m.Called(ctx, queueURL, body)
	return args.Error(0)
}

func (m *MockPublisher) GetQueueURL(ctx context.Context, queueName *string) (*string, error) {
	args := m.Called(ctx, queueName)
	url := args.String(0)
	return &url, args.Error(1)
}
