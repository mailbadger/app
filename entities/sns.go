package entities

import (
	"encoding/json"
	"time"
)

//https://docs.aws.amazon.com/ses/latest/DeveloperGuide/notification-contents.html

type Header struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type CommonHeader struct {
	From      []string  `json:"from"`
	To        []string  `json:"to"`
	Date      time.Time `json:"date"`
	Subject   string    `json:"subject"`
	MessageID string    `json:"messageId"`
}

type Mail struct {
	Timestamp        time.Time     `json:"timestamp"`
	MessageID        string        `json:"messageId"`
	Source           string        `json:"source"`
	SourceArn        string        `json:"sourceArn"`
	SourceIP         string        `json:"sourceIp"`
	SendingAccountID string        `json:"sendingAccountId"`
	Destination      []string      `json:"destination"`
	HeadersTruncated bool          `json:"headersTruncated"`
	Headers          []Header      `json:"headers"`
	CommonHeaders    *CommonHeader `json:"commonHeaders"`
}

// BouncedRecipient holds the bounced
// email address from Amazon notification.
type BouncedRecipient struct {
	EmailAddress string `json:"emailAddress"`
	Reason       string `json:"action"`
}

type ComplainedRecipient struct {
	EmailAddress string `json:"emailAddress"`
}

// Bounce field from the AWS incoming JSON notification.
type Bounce struct {
	BouncedRecipients []*BouncedRecipient `json:"bouncedRecipients"`
	BounceType        string              `json:"bounceType"`
	BounceSubType     string              `json:"bounceSubType"`
	Timestamp         time.Time           `json:"timestamp"`
	FeedbackID        string              `json:"feedbackId"`
}

type Complaint struct {
	ComplainedRecipients []*ComplainedRecipient `json:"complainedRecipients"`
	Timestamp            time.Time              `json:"timestamp"`
	FeedbackID           string                 `json:"feedbackId"`
}

type Delivery struct {
	Timestamp            time.Time `json:"timestamp"`
	ProcessingTimeMillis int64     `json:"processingTimeMillis"`
	Recipients           []string  `json:"recipients"`
	SMTPResponse         string    `json:"smtpResponse"`
	ReportingMTA         string    `json:"reportingMTA"`
	RemoteMtaIP          string    `json:"remoteMtaIp"`
}

// SNSMessage is used in the hooks action
// for processing the incoming notification messages
// with "Bounce" or "Complaint" notification type.
type SNSMessage struct {
	Type         string          `json:"Type"`
	TopicArn     string          `json:"TopicArn"`
	SubscribeURL string          `json:"SubscribeURL"`
	RawMessage   json.RawMessage `json:"Message"`
}

type SesMessage struct {
	NotificationType string     `json:"eventType"`
	Mail             Mail       `json:"mail"`
	Bounce           *Bounce    `json:"bounce"`
	Complaint        *Complaint `json:"complaint"`
	Delivery         *Delivery  `json:"delivery"`
}
