package storage

import (
	"golang.org/x/net/context"

	"github.com/FilipNikolovski/news-maily/entities"
	"github.com/FilipNikolovski/news-maily/utils/pagination"
	"github.com/gin-gonic/gin"
)

const key = "storage"

// Storage is the central interface for accessing and
// writing data in the datastore.
type Storage interface {
	GetUser(int64) (*entities.User, error)

	GetUserByAPIKey(string) (*entities.User, error)

	GetUserByUsername(string) (*entities.User, error)

	UpdateUser(*entities.User) error

	GetTemplates(int64, *pagination.Pagination)

	GetTemplate(int64, int64) (*entities.Template, error)

	GetTemplateByName(string, int64) (*entities.Template, error)

	CreateTemplate(*entities.Template) error

	UpdateTemplate(*entities.Template) error

	DeleteTemplate(int64, int64) error

	GetCampaigns(int64, *pagination.Pagination)

	GetCampaign(int64, int64) (*entities.Campaign, error)

	GetCampaignsByTemplateId(int64, int64) ([]entities.Campaign, error)

	CreateCampaign(*entities.Campaign) error

	UpdateCampaign(*entities.Campaign) error

	DeleteCampaign(int64, int64) error

	GetLists(int64, *pagination.Pagination)

	GetList(int64, int64) (*entities.List, error)

	CreateList(*entities.List) error

	UpdateList(*entities.List) error

	DeleteList(int64, int64) error

	GetSubscribers(int64, *pagination.Pagination)

	GetSubscriber(int64, int64) (*entities.Subscriber, error)

	CreateSubscriber(*entities.Subscriber) error

	UpdateSubscriber(*entities.Subscriber) error

	DeleteSubscriber(int64, int64) error
}

// SetToContext sets the storage to the context
func SetToContext(c *gin.Context, storage Storage) {
	c.Set(key, storage)
}

// GetFromContext returns the Storage associated with the context
func GetFromContext(c context.Context) Storage {
	return c.Value(key).(Storage)
}

// GetUser returns a User entity from the specified id.
func GetUser(c context.Context, id int64) (*entities.User, error) {
	return c.Value(key).(Storage).GetUser(id)
}

// GetUserByAPIKey returns a User entity from the specified api key.
func GetUserByAPIKey(c context.Context, apiKey string) (*entities.User, error) {
	return c.Value(key).(Storage).GetUserByAPIKey(apiKey)
}

// GetUserByUsername returns a User entity from the specified username.
func GetUserByUsername(c context.Context, username string) (*entities.User, error) {
	return c.Value(key).(Storage).GetUserByUsername(username)
}

// UpdateUser updates the User entity.
func UpdateUser(c context.Context, user *entities.User) error {
	return c.Value(key).(Storage).UpdateUser(user)
}

// GetTemplates populates a pagination object with a collection of
// templates by the specified user id.
func GetTemplates(c context.Context, userID int64, p *pagination.Pagination) {
	c.Value(key).(Storage).GetTemplates(userID, p)
}

// GetTemplate returns a Template entity by the given id and the user id.
func GetTemplate(c context.Context, id int64, userID int64) (*entities.Template, error) {
	return c.Value(key).(Storage).GetTemplate(id, userID)
}

// GetTemplateByName returns a Template entity by the given name and the user id.
func GetTemplateByName(c context.Context, name string, userID int64) (*entities.Template, error) {
	return c.Value(key).(Storage).GetTemplateByName(name, userID)
}

// CreateTemplate persists a new Template entity in the datastore.
func CreateTemplate(c context.Context, t *entities.Template) error {
	return c.Value(key).(Storage).CreateTemplate(t)
}

// UpdateTemplate updates the Template entity.
func UpdateTemplate(c context.Context, t *entities.Template) error {
	return c.Value(key).(Storage).UpdateTemplate(t)
}

// DeleteTemplate deletes a Template entity by the given id.
func DeleteTemplate(c context.Context, id int64, userID int64) error {
	return c.Value(key).(Storage).DeleteTemplate(id, userID)
}

// GetCampaigns populates a pagination object with a collection of
// campaigns by the specified user id.
func GetCampaigns(c context.Context, userID int64, p *pagination.Pagination) {
	c.Value(key).(Storage).GetCampaigns(userID, p)
}

// GetCampaign returns a Campaign entity by the given id and user id.
func GetCampaign(c context.Context, id int64, userID int64) (*entities.Campaign, error) {
	return c.Value(key).(Storage).GetCampaign(id, userID)
}

// GetCampaignsByTemplateId returns a Campaign entity by the given id and user id.
func GetCampaignsByTemplateId(c context.Context, templateID int64, userID int64) ([]entities.Campaign, error) {
	return c.Value(key).(Storage).GetCampaignsByTemplateId(templateID, userID)
}

// CreateCampaign persists a new Campaign entity in the datastore.
func CreateCampaign(c context.Context, campaign *entities.Campaign) error {
	return c.Value(key).(Storage).CreateCampaign(campaign)
}

// UpdateCampaign updates a Campaign entity.
func UpdateCampaign(c context.Context, campaign *entities.Campaign) error {
	return c.Value(key).(Storage).UpdateCampaign(campaign)
}

// DeleteCampaign deletes a Campaign entity by the given id.
func DeleteCampaign(c context.Context, id int64, userID int64) error {
	return c.Value(key).(Storage).DeleteCampaign(id, userID)
}

// GetLists populates a pagination object with a collection of
// lists by the specified user id.
func GetLists(c context.Context, userID int64, p *pagination.Pagination) {
	c.Value(key).(Storage).GetLists(userID, p)
}

// GetList returns a List entity by the given id and user id.
func GetList(c context.Context, id int64, userID int64) (*entities.List, error) {
	return c.Value(key).(Storage).GetList(id, userID)
}

// CreateList persists a new List entity in the datastore.
func CreateList(c context.Context, l *entities.List) error {
	return c.Value(key).(Storage).CreateList(l)
}

// UpdateList updates a List entity.
func UpdateList(c context.Context, l *entities.List) error {
	return c.Value(key).(Storage).UpdateList(l)
}

// DeleteList deletes a List entity by the given id.
func DeleteList(c context.Context, id int64, userID int64) error {
	return c.Value(key).(Storage).DeleteList(id, userID)
}

// GetSubscribers populates a pagination object with a collection of
// subscribers by the specified user id.
func GetSubscribers(c context.Context, userID int64, p *pagination.Pagination) {
	c.Value(key).(Storage).GetSubscribers(userID, p)
}

// GetSubscriber returns a Subscriber entity by the given id and user id.
func GetSubscriber(c context.Context, id int64, userID int64) (*entities.Subscriber, error) {
	return c.Value(key).(Storage).GetSubscriber(id, userID)
}

// CreateSubscriber persists a new Subscriber entity in the datastore.
func CreateSubscriber(c context.Context, s *entities.Subscriber) error {
	return c.Value(key).(Storage).CreateSubscriber(s)
}

// UpdateSubscriber updates a Subscriber entity.
func UpdateSubscriber(c context.Context, s *entities.Subscriber) error {
	return c.Value(key).(Storage).UpdateSubscriber(s)
}

// DeleteSubscriber deletes a Subscriber entity by the given id.
func DeleteSubscriber(c context.Context, id int64, userID int64) error {
	return c.Value(key).(Storage).DeleteSubscriber(id, userID)
}
