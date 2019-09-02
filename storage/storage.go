package storage

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/news-maily/app/entities"
	"github.com/news-maily/app/utils/pagination"
)

const key = "storage"

// Storage is the central interface for accessing and
// writing data in the datastore.
type Storage interface {
	GetUser(int64) (*entities.User, error)
	GetUserByUUID(string) (*entities.User, error)
	GetUserByUsername(string) (*entities.User, error)
	GetActiveUserByUsername(string) (*entities.User, error)
	CreateUser(*entities.User) error
	UpdateUser(*entities.User) error

	GetSession(sessionID string) (*entities.Session, error)
	CreateSession(s *entities.Session) error
	DeleteSession(userID int64) error

	GetCampaigns(int64, *pagination.Cursor)
	GetCampaign(int64, int64) (*entities.Campaign, error)
	GetCampaignByName(name string, userID int64) (*entities.Campaign, error)
	GetCampaignsByTemplateName(string, int64) ([]entities.Campaign, error)
	CreateCampaign(*entities.Campaign) error
	UpdateCampaign(*entities.Campaign) error
	DeleteCampaign(int64, int64) error

	GetSegments(int64, *pagination.Cursor)
	GetSegmentsByIDs(userID int64, ids []int64) ([]entities.Segment, error)
	GetSegment(int64, int64) (*entities.Segment, error)
	GetSegmentByName(name string, userID int64) (*entities.Segment, error)
	CreateSegment(*entities.Segment) error
	UpdateSegment(*entities.Segment) error
	DeleteSegment(int64, int64) error
	AppendSubscribers(*entities.Segment) error
	DetachSubscribers(*entities.Segment) error

	GetSubscribers(int64, *pagination.Cursor)
	GetSubscribersBySegmentID(int64, int64, *pagination.Cursor)
	GetSubscriber(int64, int64) (*entities.Subscriber, error)
	GetSubscribersByIDs([]int64, int64) ([]entities.Subscriber, error)
	GetSubscriberByEmail(string, int64) (*entities.Subscriber, error)
	GetDistinctSubscribersBySegmentIDs(
		listIDs []int64,
		userID int64,
		blacklisted, active bool,
		timestamp time.Time,
		nextID, limit int64,
	) ([]entities.Subscriber, error)
	CreateSubscriber(*entities.Subscriber) error
	UpdateSubscriber(*entities.Subscriber) error
	BlacklistSubscriber(userID int64, email string) error
	DeleteSubscriber(int64, int64) error

	GetAPIKeys(userID int64) []*entities.APIKey
	GetAPIKey(identifier string) (*entities.APIKey, error)
	CreateAPIKey(ak *entities.APIKey) error
	UpdateAPIKey(ak *entities.APIKey) error
	DeleteAPIKey(id, userID int64) error

	GetSesKeys(userID int64) (*entities.SesKeys, error)
	CreateSesKeys(s *entities.SesKeys) error
	DeleteSesKeys(userID int64) error

	CreateSendBulkLog(l *entities.SendBulkLog) error
	CountLogsByUUID(uuid string) (int, error)

	CreateBounce(b *entities.Bounce) error
	CreateComplaint(c *entities.Complaint) error
	CreateSend(s *entities.Send) error
	CreateClick(c *entities.Click) error
	CreateOpen(o *entities.Open) error
	CreateDelivery(d *entities.Delivery) error
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
	return GetFromContext(c).GetUser(id)
}

// GetUserByUUID returns a User entity from the specified uuid.
func GetUserByUUID(c context.Context, uuid string) (*entities.User, error) {
	return GetFromContext(c).GetUserByUUID(uuid)
}

// GetUserByUsername returns a User entity from the specified username.
func GetUserByUsername(c context.Context, username string) (*entities.User, error) {
	return GetFromContext(c).GetUserByUsername(username)
}

// GetActiveUserByUsername returns an active User entity from the specified username.
func GetActiveUserByUsername(c context.Context, username string) (*entities.User, error) {
	return GetFromContext(c).GetActiveUserByUsername(username)
}

// CreateUser persists a new User entity in the datastore.
func CreateUser(c context.Context, user *entities.User) error {
	return GetFromContext(c).CreateUser(user)
}

// UpdateUser updates the User entity.
func UpdateUser(c context.Context, user *entities.User) error {
	return GetFromContext(c).UpdateUser(user)
}

// GetSession returns the session by the given session id.
func GetSession(c context.Context, sessionID string) (*entities.Session, error) {
	return GetFromContext(c).GetSession(sessionID)
}

// CreateSession adds a new session in the database.
func CreateSession(c context.Context, s *entities.Session) error {
	return GetFromContext(c).CreateSession(s)
}

// DeleteSession deletes a session by the given user id from the database.
func DeleteSession(c context.Context, userID int64) error {
	return GetFromContext(c).DeleteSession(userID)
}

// GetCampaigns populates a pagination object with a collection of
// campaigns by the specified user id.
func GetCampaigns(c context.Context, userID int64, p *pagination.Cursor) {
	GetFromContext(c).GetCampaigns(userID, p)
}

// GetCampaign returns a Campaign entity by the given id and user id.
func GetCampaign(c context.Context, id, userID int64) (*entities.Campaign, error) {
	return GetFromContext(c).GetCampaign(id, userID)
}

// GetCampaignByName returns a Campaign entity by the given name and user id.
func GetCampaignByName(c context.Context, name string, userID int64) (*entities.Campaign, error) {
	return GetFromContext(c).GetCampaignByName(name, userID)
}

// GetCampaignsByTemplateName returns a collection of campaigns by the given template name and user id.
func GetCampaignsByTemplateName(c context.Context, templateName string, userID int64) ([]entities.Campaign, error) {
	return GetFromContext(c).GetCampaignsByTemplateName(templateName, userID)
}

// CreateCampaign persists a new Campaign entity in the datastore.
func CreateCampaign(c context.Context, campaign *entities.Campaign) error {
	return GetFromContext(c).CreateCampaign(campaign)
}

// UpdateCampaign updates a Campaign entity.
func UpdateCampaign(c context.Context, campaign *entities.Campaign) error {
	return GetFromContext(c).UpdateCampaign(campaign)
}

// DeleteCampaign deletes a Campaign entity by the given id.
func DeleteCampaign(c context.Context, id, userID int64) error {
	return GetFromContext(c).DeleteCampaign(id, userID)
}

// GetSegments populates a pagination object with a collection of
// lists by the specified user id.
func GetSegments(c context.Context, userID int64, p *pagination.Cursor) {
	GetFromContext(c).GetSegments(userID, p)
}

// GetSegmentsByIDs fetches lists by user id and the given ids
func GetSegmentsByIDs(c context.Context, userID int64, ids []int64) ([]entities.Segment, error) {
	return GetFromContext(c).GetSegmentsByIDs(userID, ids)
}

// GetSegment returns a Segment entity by the given id and user id.
func GetSegment(c context.Context, id, userID int64) (*entities.Segment, error) {
	return GetFromContext(c).GetSegment(id, userID)
}

// GetSegmentByName returns a Campaign entity by the given name and user id.
func GetSegmentByName(c context.Context, name string, userID int64) (*entities.Segment, error) {
	return GetFromContext(c).GetSegmentByName(name, userID)
}

// CreateSegment persists a new Segment entity in the datastore.
func CreateSegment(c context.Context, l *entities.Segment) error {
	return GetFromContext(c).CreateSegment(l)
}

// UpdateSegment updates a Segment entity.
func UpdateSegment(c context.Context, l *entities.Segment) error {
	return GetFromContext(c).UpdateSegment(l)
}

// DeleteSegment deletes a Segment entity by the given id.
func DeleteSegment(c context.Context, id, userID int64) error {
	return GetFromContext(c).DeleteSegment(id, userID)
}

// AppendSubscribers appends subscribers to the existing association.
func AppendSubscribers(c context.Context, l *entities.Segment) error {
	return GetFromContext(c).AppendSubscribers(l)
}

// DetachSubscribers deletes subscribers from the list.
func DetachSubscribers(c context.Context, l *entities.Segment) error {
	return GetFromContext(c).DetachSubscribers(l)
}

// GetSubscribers populates a pagination object with a collection of
// subscribers by the specified user id.
func GetSubscribers(c context.Context, userID int64, p *pagination.Cursor) {
	GetFromContext(c).GetSubscribers(userID, p)
}

// GetSubscribersBySegmentID populates a pagination object with a collection of
// subscribers by the specified user id and list id.
func GetSubscribersBySegmentID(c context.Context, segmentID, userID int64, p *pagination.Cursor) {
	GetFromContext(c).GetSubscribersBySegmentID(segmentID, userID, p)
}

// GetSubscriber returns a Subscriber entity by the given id and user id.
func GetSubscriber(c context.Context, id, userID int64) (*entities.Subscriber, error) {
	return GetFromContext(c).GetSubscriber(id, userID)
}

// GetSubscribersByIDs returns a Subscriber entity by the given ids and user id.
func GetSubscribersByIDs(c context.Context, ids []int64, userID int64) ([]entities.Subscriber, error) {
	return GetFromContext(c).GetSubscribersByIDs(ids, userID)
}

// GetSubscriberByEmail returns a Subscriber entity by the given email and user id.
func GetSubscriberByEmail(c context.Context, email string, userID int64) (*entities.Subscriber, error) {
	return GetFromContext(c).GetSubscriberByEmail(email, userID)
}

// GetDistinctSubscribersBySegmentIDs fetches all distinct subscribers by user id and list ids
func GetDistinctSubscribersBySegmentIDs(
	c context.Context,
	listIDs []int64,
	userID int64,
	blacklisted, active bool,
	timestamp time.Time,
	nextID, limit int64,
) ([]entities.Subscriber, error) {
	return GetFromContext(c).GetDistinctSubscribersBySegmentIDs(listIDs, userID, blacklisted, active, timestamp, nextID, limit)
}

// CreateSubscriber persists a new Subscriber entity in the datastore.
func CreateSubscriber(c context.Context, s *entities.Subscriber) error {
	return GetFromContext(c).CreateSubscriber(s)
}

// UpdateSubscriber updates a Subscriber entity.
func UpdateSubscriber(c context.Context, s *entities.Subscriber) error {
	return GetFromContext(c).UpdateSubscriber(s)
}

// BlacklistSubscriber blacklists a Subscriber entity by the given email.
func BlacklistSubscriber(c context.Context, userID int64, email string) error {
	return GetFromContext(c).BlacklistSubscriber(userID, email)
}

// DeleteSubscriber deletes a Subscriber entity by the given id.
func DeleteSubscriber(c context.Context, id, userID int64) error {
	return GetFromContext(c).DeleteSubscriber(id, userID)
}

// GetAPIKeys returns a list of APIKey entities.
func GetAPIKeys(c context.Context, userID int64) []*entities.APIKey {
	return GetFromContext(c).GetAPIKeys(userID)
}

// GetAPIKey returns an APIKey entity by the given identifier.
func GetAPIKey(c context.Context, identifier string) (*entities.APIKey, error) {
	return GetFromContext(c).GetAPIKey(identifier)
}

// CreateAPIKey persists a new APIKey entity in the datastore.
func CreateAPIKey(c context.Context, ak *entities.APIKey) error {
	return GetFromContext(c).CreateAPIKey(ak)
}

// UpdateAPIKey updates an APIKey entity.
func UpdateAPIKey(c context.Context, ak *entities.APIKey) error {
	return GetFromContext(c).UpdateAPIKey(ak)
}

// DeleteAPIKey deletes an APIKey entity by the given id.
func DeleteAPIKey(c context.Context, id, userID int64) error {
	return GetFromContext(c).DeleteAPIKey(id, userID)
}

// GetSesKeys returns the SES keys by the given user id
func GetSesKeys(c context.Context, userID int64) (*entities.SesKeys, error) {
	return GetFromContext(c).GetSesKeys(userID)
}

// CreateSesKeys adds new SES keys in the database.
func CreateSesKeys(c context.Context, s *entities.SesKeys) error {
	return GetFromContext(c).CreateSesKeys(s)
}

// DeleteSesKeys deletes SES keys configuration by the given user ID.
func DeleteSesKeys(c context.Context, userID int64) error {
	return GetFromContext(c).DeleteSesKeys(userID)
}

// CreateBounce adds new bounce in the database.
func CreateBounce(c context.Context, b *entities.Bounce) error {
	return GetFromContext(c).CreateBounce(b)
}

// CreateComplaint adds new complaint in the database.
func CreateComplaint(c context.Context, compl *entities.Complaint) error {
	return GetFromContext(c).CreateComplaint(compl)
}

// CreateSend adds new send in the database.
func CreateSend(c context.Context, send *entities.Send) error {
	return GetFromContext(c).CreateSend(send)
}

// CreateClick adds new click in the database.
func CreateClick(c context.Context, click *entities.Click) error {
	return GetFromContext(c).CreateClick(click)
}

// CreateOpen adds new open in the database.
func CreateOpen(c context.Context, open *entities.Open) error {
	return GetFromContext(c).CreateOpen(open)
}

// CreateDelivery adds new delivery in the database.
func CreateDelivery(c context.Context, d *entities.Delivery) error {
	return GetFromContext(c).CreateDelivery(d)
}
