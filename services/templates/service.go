package templates

import (
	"context"

	"github.com/mailbadger/app/storage"
)

type Service interface {
	ListTemplates(c context.Context, userID int64, p *storage.PaginationCursor, scopeMap map[string]string) error
}

type service struct {
}

func New() Service {
	return &service{}
}

// ListTemplates populates a pagination object with a collection of
// templates by the specified user id.
func (s service) ListTemplates(c context.Context, userID int64, p *storage.PaginationCursor, scopeMap map[string]string) error {
	return storage.ListTemplates(c, userID, p, scopeMap)
}
