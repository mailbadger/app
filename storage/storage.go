package storage

import (
	"golang.org/x/net/context"

	"github.com/FilipNikolovski/news-maily/entities"
	"github.com/FilipNikolovski/news-maily/routes/middleware"
)

const key = "storage"

type Storage interface {
	GetUser(id int64) (entities.User, error)

	GetUserByApiKey(api_key string) (entities.User, error)

	UpdateUser(user *entities.User) error

	GetTemplates(user_id int64, p *middleware.Pagination)

	GetTemplate(id int64, user_id int64) (entities.Template, error)

	CreateTemplate(t *entities.Template) error

	UpdateTemplate(t *entities.Template) error

	DeleteTemplate(id int64, user_id int64) error
}

func GetUser(c context.Context, id int64) (entities.User, error) {
	return c.Value(key).(Storage).GetUser(id)
}

func GetUserByApiKey(c context.Context, api_key string) (entities.User, error) {
	return c.Value(key).(Storage).GetUserByApiKey(api_key)
}

func UpdateUser(c context.Context, user *entities.User) error {
	return c.Value(key).(Storage).UpdateUser(user)
}

// GetTemplates populates the Pagination object with a collection of templates
// and page data.
func GetTemplates(c context.Context, user_id int64, p *middleware.Pagination) {
	c.Value(key).(Storage).GetTemplates(user_id, p)
}

func GetTemplate(c context.Context, id int64, user_id int64) (entities.Template, error) {
	return c.Value(key).(Storage).GetTemplate(id, user_id)
}

func CreateTemplate(c context.Context, t *entities.Template) error {
	return c.Value(key).(Storage).CreateTemplate(t)
}

func UpdateTemplate(c context.Context, t *entities.Template) error {
	return c.Value(key).(Storage).UpdateTemplate(t)
}

func DeleteTemplate(c context.Context, id int64, user_id int64) error {
	return c.Value(key).(Storage).DeleteTemplate(id, user_id)
}
