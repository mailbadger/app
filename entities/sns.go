package entities

import (
	"time"
)

// https://docs.aws.amazon.com/ses/latest/DeveloperGuide/notification-contents.html

type Header struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type CommonHeader struct {
	From      []string `json:"from"`
	To        []string `json:"to"`
	Date      string   `json:"date"`
	Subject   string   `json:"subject"`
	MessageID string   `json:"messageId"`
}

// Mail field from the AWS incoming JSON notification.
type Mail struct {
	Timestamp        time.Time           `json:"timestamp"`
	MessageID        string              `json:"messageId"`
	Source           string              `json:"source"`
	SourceArn        string              `json:"sourceArn"`
	SourceIP         string              `json:"sourceIp"`
	SendingAccountID string              `json:"sendingAccountId"`
	Destination      []string            `json:"destination"`
	HeadersTruncated bool                `json:"headersTruncated"`
	Headers          []Header            `json:"headers"`
	CommonHeaders    *CommonHeader       `json:"commonHeaders"`
	Tags             map[string][]string `json:"tags"`
}

// BouncedRecipient holds the bounced
// email address from Amazon notification.
type BouncedRecipient struct {
	EmailAddress   string `json:"emailAddress"`
	Action         string `json:"action"`
	Status         string `json:"status"`
	DiagnosticCode string `json:"diagnosticCode"`
}

// BounceMsg field from the AWS incoming JSON notification.
type BounceMsg struct {
	BouncedRecipients []*BouncedRecipient `json:"bouncedRecipients"`
	BounceType        string              `json:"bounceType"`
	BounceSubType     string              `json:"bounceSubType"`
	Timestamp         time.Time           `json:"timestamp"`
	FeedbackID        string              `json:"feedbackId"`
	ReportingMTA      string              `json:"reportingMTA"`
	RemoteMTAIp       string              `json:"remoteMtaIp"`
}

type ComplainedRecipient struct {
	EmailAddress string `json:"emailAddress"`
}

// ComplaintMsg field from the AWS incoming JSON notification.
type ComplaintMsg struct {
	ComplainedRecipients  []*ComplainedRecipient `json:"complainedRecipients"`
	Timestamp             time.Time              `json:"timestamp"`
	FeedbackID            string                 `json:"feedbackId"`
	UserAgent             string                 `json:"userAgent"`
	ComplaintFeedbackType string                 `json:"complaintFeedbackType"`
}

// DeliveryMsg field from the AWS incoming JSON notification.
type DeliveryMsg struct {
	Timestamp            time.Time `json:"timestamp"`
	ProcessingTimeMillis int64     `json:"processingTimeMillis"`
	Recipients           []string  `json:"recipients"`
	SMTPResponse         string    `json:"smtpResponse"`
	ReportingMTA         string    `json:"reportingMTA"`
	RemoteMtaIP          string    `json:"remoteMtaIp"`
}

// RenderingFailureMsg field from the AWS incoming JSON notification.
type RenderingFailureMsg struct {
	ErrorMessage string `json:"errorMessage"`
	TemplateName string `json:"templateName"`
}

// ClickMsg field from the AWS incoming JSON notification.
type ClickMsg struct {
	Timestamp time.Time           `json:"timestamp"`
	IPAddress string              `json:"ipAddress"`
	UserAgent string              `json:"userAgent"`
	Link      string              `json:"link"`
	LinkTags  map[string][]string `json:"linkTags"`
}

// ClickMsg field from the AWS incoming JSON notification.
type OpenMsg struct {
	Timestamp time.Time `json:"timestamp"`
	IPAddress string    `json:"ipAddress"`
	UserAgent string    `json:"userAgent"`
}

// SesMessage represents the message that is sent by the SNS topic.
type SesMessage struct {
	NotificationType string               `json:"eventType"`
	Mail             Mail                 `json:"mail"`
	Bounce           *BounceMsg           `json:"bounce"`
	Complaint        *ComplaintMsg        `json:"complaint"`
	Delivery         *DeliveryMsg         `json:"delivery"`
	RenderingFailure *RenderingFailureMsg `json:"failure"`
	Click            *ClickMsg            `json:"click"`
	Open             *OpenMsg             `json:"open"`
}
