package emails

import (
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/stretchr/testify/mock"
)

// MockSender structure with ses api mock for testing
type MockSender struct {
	mock.Mock
}

func (m *MockSender) SendEmail(input *ses.SendEmailInput) (*ses.SendEmailOutput, error) {
	args := m.Called(input)
	return nil, args.Error(1)
}

func (m *MockSender) CreateConfigurationSet(input *ses.CreateConfigurationSetInput) (*ses.CreateConfigurationSetOutput, error) {
	args := m.Called(input)
	return nil, args.Error(1)
}

func (m *MockSender) DescribeConfigurationSet(input *ses.DescribeConfigurationSetInput) (*ses.DescribeConfigurationSetOutput, error) {
	args := m.Called(input)
	return nil, args.Error(1)
}

func (m *MockSender) CreateConfigurationSetEventDestination(
	input *ses.CreateConfigurationSetEventDestinationInput,
) (*ses.CreateConfigurationSetEventDestinationOutput, error) {
	args := m.Called(input)
	return nil, args.Error(1)
}

func (m *MockSender) DeleteConfigurationSet(input *ses.DeleteConfigurationSetInput) (*ses.DeleteConfigurationSetOutput, error) {
	args := m.Called(input)
	return nil, args.Error(1)
}

func (m *MockSender) GetSendQuota(input *ses.GetSendQuotaInput) (*ses.GetSendQuotaOutput, error) {
	args := m.Called(input)
	return nil, args.Error(1)
}
