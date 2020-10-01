package params

import "strings"

// PostSubscriber represents request body for POST /api/subscribers
type PostSubscriber struct {
	Name       string            `form:"name" validate:"omitempty,min=1,max=191"`
	Email      string            `form:"email" validate:"required,email"`
	SegmentIDs []int64           `form:"segments[]" validate:"omitempty"`
	Metadata   map[string]string `form:"metadata" validate:"omitempty,dive,alphanumhyphen"`
}

func (p *PostSubscriber) TrimSpaces() {
	p.Name = strings.TrimSpace(p.Name)
}

// PutSubscriber represents request body for PUT /api/subscribers/:id
type PutSubscriber struct {
	Name       string            `form:"name" validate:"omitempty,min=1,max=191"`
	SegmentIDs []int64           `form:"segments[]" validate:"omitempty"`
	Metadata   map[string]string `form:"metadata" validate:"omitempty,dive,alphanumhyphen"`
}

func (p *PutSubscriber) TrimSpaces() {
	p.Name = strings.TrimSpace(p.Name)
}

// PostUnsubscribe represents request body for POST /api/unsubscribe
type PostUnsubscribe struct {
	Email string `form:"email" validate:"required,email"`
	UUID  string `form:"uuid" validate:"required,uuid"`
	Token string `form:"t" validate:"required"`
}

func (p *PostUnsubscribe) TrimSpaces() {
	p.Email = strings.TrimSpace(p.Email)
	p.UUID = strings.TrimSpace(p.UUID)
	p.Token = strings.TrimSpace(p.Token)
}

// ImportSubscribers represents request body for POST /api/subscribers/import
type ImportSubscribers struct {
	Filename   string  `form:"filename" validate:"required"`
	SegmentIDs []int64 `form:"segments[]" validate:"omitempty"`
}

func (p *ImportSubscribers) TrimSpaces() {
	p.Filename = strings.TrimSpace(p.Filename)
}

// BulkRemoveSubscribers represents request body for POST /api/subscribers/bulk-remove
type BulkRemoveSubscribers struct {
	Filename string `form:"filename" validate:"required"`
}

func (p *BulkRemoveSubscribers) TrimSpaces() {
	p.Filename = strings.TrimSpace(p.Filename)
}
