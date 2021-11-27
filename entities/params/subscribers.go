package params

import "strings"

// PostSubscriber represents request body for POST /api/subscribers
type PostSubscriber struct {
	Name       string            `json:"name" validate:"omitempty,min=1,max=191"`
	Email      string            `json:"email" validate:"required,email"`
	SegmentIDs []int64           `json:"segments" validate:"omitempty"`
	Metadata   map[string]string `json:"metadata" validate:"omitempty,dive,keys,required,alphanumhyphen,endkeys,required"`
}

func (p *PostSubscriber) TrimSpaces() {
	p.Name = strings.TrimSpace(p.Name)
}

// PutSubscriber represents request body for PUT /api/subscribers/:id
type PutSubscriber struct {
	Name       string            `json:"name" validate:"omitempty,min=1,max=191"`
	SegmentIDs []int64           `json:"segments" validate:"omitempty"`
	Metadata   map[string]string `json:"metadata" validate:"omitempty,dive,keys,required,alphanumhyphen,endkeys,required"`
}

func (p *PutSubscriber) TrimSpaces() {
	p.Name = strings.TrimSpace(p.Name)
}

// PostUnsubscribe represents request body for POST /api/unsubscribe
type PostUnsubscribe struct {
	Email string `json:"email" validate:"required,email"`
	UUID  string `json:"uuid" validate:"required,uuid"`
	Token string `json:"t" validate:"required"`
}

func (p *PostUnsubscribe) TrimSpaces() {
	p.Email = strings.TrimSpace(p.Email)
	p.UUID = strings.TrimSpace(p.UUID)
	p.Token = strings.TrimSpace(p.Token)
}

// ImportSubscribers represents request body for POST /api/subscribers/import
type ImportSubscribers struct {
	Filename   string  `json:"filename" validate:"required"`
	SegmentIDs []int64 `json:"segments" validate:"omitempty"`
}

func (p *ImportSubscribers) TrimSpaces() {
	p.Filename = strings.TrimSpace(p.Filename)
}

// BulkRemoveSubscribers represents request body for POST /api/subscribers/bulk-remove
type BulkRemoveSubscribers struct {
	Filename string `json:"filename" validate:"required"`
}

func (p *BulkRemoveSubscribers) TrimSpaces() {
	p.Filename = strings.TrimSpace(p.Filename)
}
