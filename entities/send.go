package entities

type Send struct {
	Id               int64  `json:"id" gorm:"column:id; primary_key:yes"`
	UserId           int64  `json:"-" gorm:"column:user_id; index"`
	CampaignId       int64  `json:"campaign_id"`
	MessageID        string `json:"message_id"`
	Source           string `json:"source"`
	SourceArn        string `json:"source_arn"`
	SourceIP         string `json:"source_ip"`
	SendingAccountID string `json:"sending_account_id"`
	Destination      string `json:"destination"`
}
