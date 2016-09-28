package storage

import (
	"golang.org/x/net/context"

	"github.com/FilipNikolovski/news-maily/entities"
	"github.com/FilipNikolovski/news-maily/utils/pagination"
	"github.com/gin-gonic/gin"
)

const key = "storage"

type Storage interface {
	GetUser(int64) (*entities.User, error)

	GetUserByApiKey(string) (*entities.User, error)

	UpdateUser(*entities.User) error

	GetTemplates(int64, *pagination.Pagination)

	GetTemplate(int64, int64) (*entities.Template, error)

	CreateTemplate(*entities.Template) error

	UpdateTemplate(*entities.Template) error

	DeleteTemplate(int64, int64) error
}

// SetToContext sets the storage to the context
func SetToContext(c *gin.Context, storage Storage) {
	c.Set(key, storage)
}

// GetFromContext returns the Storage associated with the context
func GetFromContext(c context.Context) Storage {
	return c.Value(key).(Storage)
}

func GetUser(c context.Context, id int64) (*entities.User, error) {
	return c.Value(key).(Storage).GetUser(id)
}

func GetUserByApiKey(c context.Context, api_key string) (*entities.User, error) {
	return c.Value(key).(Storage).GetUserByApiKey(api_key)
}

func UpdateUser(c context.Context, user *entities.User) error {
	return c.Value(key).(Storage).UpdateUser(user)
}

// GetTemplates populates the Pagination object with a collection of templates
// and page data.
func GetTemplates(c context.Context, user_id int64, p *pagination.Pagination) {
	c.Value(key).(Storage).GetTemplates(user_id, p)
}

func GetTemplate(c context.Context, id int64, user_id int64) (*entities.Template, error) {
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
