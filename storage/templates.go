package storage

import (
	"github.com/mailbadger/app/entities"
)

func (db *store) ListTemplates(userID int64, p *PaginationCursor, scopeMap map[string]string) error {
	p.SetCollection(&[]entities.TemplatesCollection{})
	p.SetResource("templates")

	for k, v := range scopeMap {
		if k == "name" {
			p.AddScope(NameLike(v))
		}
	}

	p.SetQuery(db.Table(p.Resource).
		Where("user_id = ?", userID).
		Order("created_at desc, id desc").
		Limit(p.PerPage))

	return db.Paginate(p, userID)
}
