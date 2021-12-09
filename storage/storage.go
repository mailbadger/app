package storage

import (
	"context"
	"time"
	
	"github.com/gin-gonic/gin"
	
	"github.com/mailbadger/app/entities"
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
	DeleteUser(user *entities.User) error
	
	GetBoundariesByType(t string) (*entities.Boundaries, error)
	GetRole(name string) (*entities.Role, error)
	
	GetSession(sessionID string) (*entities.Session, error)
	CreateSession(s *entities.Session) error
	DeleteSession(sessionID string) error
	DeleteAllSessionsForUser(userID int64) error
	
	GetCampaigns(int64, *PaginationCursor, map[string]string) error
	GetCampaign(int64, int64) (*entities.Campaign, error)
	GetCampaignByName(name string, userID int64) (*entities.Campaign, error)
	CreateCampaign(*entities.Campaign) error
	UpdateCampaign(*entities.Campaign) error
	DeleteCampaign(int64, int64) error
	GetMonthlyTotalCampaigns(userID int64) (int64, error)
	GetCampaignOpens(campaignID, userID int64, p *PaginationCursor) error
	GetClicksStats(campaignID, userID int64) (*entities.ClicksStats, error)
	GetOpensStats(campaignID, userID int64) (*entities.OpensStats, error)
	GetTotalSends(campaignID, userID int64) (int64, error)
	GetTotalDelivered(campaignID, userID int64) (int64, error)
	GetTotalBounces(campaignID, userID int64) (int64, error)
	GetTotalComplaints(campaignID, userID int64) (int64, error)
	GetCampaignClicksStats(int64, int64) ([]entities.ClicksStats, error)
	GetCampaignComplaints(campaignID, userID int64, p *PaginationCursor) error
	GetCampaignBounces(campaignID, userID int64, p *PaginationCursor) error
	DeleteAllCampaignsForUser(userID int64) error
	LogFailedCampaign(c *entities.Campaign, description string) error
	DeleteAllCampaignFailedLogsForUser(userID int64) error
	
	CreateCampaignSchedule(c *entities.CampaignSchedule) error
	DeleteCampaignSchedule(campaignID int64) error
	GetScheduledCampaigns(time time.Time) ([]entities.CampaignSchedule, error)
	
	GetSegments(int64, *PaginationCursor) error
	GetSegmentsByIDs(userID int64, ids []int64) ([]entities.Segment, error)
	GetSegment(int64, int64) (*entities.Segment, error)
	GetSegmentByName(name string, userID int64) (*entities.Segment, error)
	GetTotalSegments(userID int64) (int64, error)
	CreateSegment(*entities.Segment) error
	UpdateSegment(*entities.Segment) error
	DeleteSegment(int64, int64) error
	AppendSubscribers(*entities.Segment) error
	DetachSubscribers(*entities.Segment) error
	DeleteAllSegmentsForUser(userID int64) error
	
	GetSubscribers(userID int64, p *PaginationCursor, scopeMap map[string]string) error
	GetSubscribersBySegmentID(int64, int64, *PaginationCursor) error
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
	DeactivateSubscriber(userID int64, email string) error
	DeleteSubscriber(int64, int64) error
	DeleteSubscriberByEmail(string, int64) error
	GetTotalSubscribers(int64) (int64, error)
	GetTotalSubscribersBySegment(segmentID, userID int64) (int64, error)
	SeekSubscribersByUserID(userID int64, nextID int64, limit int64) ([]entities.Subscriber, error)
	GetAllSubscribersForUser(userID int64) ([]entities.Subscriber, error)
	
	GetAPIKeys(userID int64) ([]*entities.APIKey, error)
	GetAPIKey(identifier string) (*entities.APIKey, error)
	CreateAPIKey(ak *entities.APIKey) error
	UpdateAPIKey(ak *entities.APIKey) error
	DeleteAPIKey(id, userID int64) error
	
	GetSesKeys(userID int64) (*entities.SesKeys, error)
	CreateSesKeys(s *entities.SesKeys) error
	DeleteSesKeys(userID int64) error
	
	GetToken(token string) (*entities.Token, error)
	CreateToken(s *entities.Token) error
	DeleteToken(token string) error
	DeleteAllTokensForUser(userID int64) error
	
	CreateSendLog(l *entities.SendLog) error
	CountLogsByUUID(id string) (int, error)
	CountLogsByStatus(status string) (int, error)
	GetSendLogByUUID(id string) (*entities.SendLog, error)
	DeleteAllSendLogsForUser(userID int64) error
	
	CreateBounce(b *entities.Bounce) error
	DeleteAllBouncesForUser(userID int64) error
	CreateComplaint(c *entities.Complaint) error
	DeleteAllComplaintsForUser(userID int64) error
	CreateSend(s *entities.Send) error
	DeleteAllSendsForUser(userID int64) error
	CreateClick(c *entities.Click) error
	DeleteAllClicksForUser(userID int64) error
	CreateOpen(o *entities.Open) error
	DeleteAllOpensForUser(userID int64) error
	CreateDelivery(d *entities.Delivery) error
	DeleteAllDeliveriesForUser(userID int64) error
	
	CreateReport(r *entities.Report) error
	UpdateReport(r *entities.Report) error
	GetReportByFilename(filename string, userID int64) (*entities.Report, error)
	GetRunningReportForUser(userID int64) (*entities.Report, error)
	GetNumberOfReportsForDate(userID int64, time time.Time) (int64, error)
	DeleteAllReportsForUser(userID int64) error
	
	CreateTemplate(t *entities.Template) error
	UpdateTemplate(t *entities.Template) error
	GetTemplateByName(name string, userID int64) (*entities.Template, error)
	GetTemplate(templateID int64, userID int64) (*entities.Template, error)
	GetTemplates(userID int64, p *PaginationCursor, scopeMap map[string]string) error
	DeleteTemplate(templateID int64, userID int64) error
	GetAllTemplatesForUser(userID int64) ([]entities.Template, error)
	
	DeleteAllEventsForUser(userID int64) error
	
	GetJobByName(name string) (*entities.Job, error)
	UpdateJob(job *entities.Job) error
	CreateSubscriberMetrics(sm []*entities.SubscribersMetrics, job *entities.Job) error
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
func DeleteSession(c context.Context, sessionID string) error {
	return GetFromContext(c).DeleteSession(sessionID)
}

func GetBoundariesByType(c context.Context, t string) (*entities.Boundaries, error) {
	return GetFromContext(c).GetBoundariesByType(t)
}

func GetRole(c context.Context, name string) (*entities.Role, error) {
	return GetFromContext(c).GetRole(name)
}

// GetCampaigns populates a pagination object with a collection of
// campaigns by the specified user id.
func GetCampaigns(c context.Context, userID int64, p *PaginationCursor, scopeMap map[string]string) error {
	return GetFromContext(c).GetCampaigns(userID, p, scopeMap)
}

// GetCampaign returns a Campaign entity by the given id and user id.
func GetCampaign(c context.Context, id, userID int64) (*entities.Campaign, error) {
	return GetFromContext(c).GetCampaign(id, userID)
}

// GetCampaignByName returns a Campaign entity by the given name and user id.
func GetCampaignByName(c context.Context, name string, userID int64) (*entities.Campaign, error) {
	return GetFromContext(c).GetCampaignByName(name, userID)
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

// GetCampaignBounces populates a pagination object with a collection of
// bounce by the specified campaign id
func GetCampaignBounces(c context.Context, campaignID, userID int64, p *PaginationCursor) error {
	return GetFromContext(c).GetCampaignBounces(campaignID, userID, p)
}

// GetCampaignComplaints populates a pagination object with a collection of
// complaints by the specified campaign id
func GetCampaignComplaints(c context.Context, campaignID, userID int64, p *PaginationCursor) error {
	return GetFromContext(c).GetCampaignComplaints(campaignID, userID, p)
}

// GetClicksStats populates a ClickStats object with a results of
// clicks by the specified campaign id
func GetClicksStats(c context.Context, campaignID, userID int64) (*entities.ClicksStats, error) {
	return GetFromContext(c).GetClicksStats(campaignID, userID)
}

// GetOpensStats populates a ClickStats object with a results of
// opens by the specified campaign id
func GetOpensStats(c context.Context, campaignID, userID int64) (*entities.OpensStats, error) {
	return GetFromContext(c).GetOpensStats(campaignID, userID)
}

// GetCampaignOpens populates a pagination object with a collection of
// open by the specified campaign id
func GetCampaignOpens(c context.Context, campaignID, userID int64, p *PaginationCursor) error {
	return GetFromContext(c).GetCampaignOpens(campaignID, userID, p)
}

// GetMonthlyTotalCampaigns returns the total number of campaigns for a specified user in the current month
func GetMonthlyTotalCampaigns(c context.Context, userID int64) (int64, error) {
	return GetFromContext(c).GetMonthlyTotalCampaigns(userID)
}

// LogFailedCampaign updates campaign status to failed & stores campaign failed log record.
func LogFailedCampaign(c context.Context, ca *entities.Campaign, description string) error {
	return GetFromContext(c).LogFailedCampaign(ca, description)
}

// CreateCampaignSchedule creates new schedule for campaign.
func CreateCampaignSchedule(c context.Context, sc *entities.CampaignSchedule) error {
	return GetFromContext(c).CreateCampaignSchedule(sc)
}

// DeleteCampaignSchedule deletes campaign schedule.
func DeleteCampaignSchedule(c context.Context, campaignID int64) error {
	return GetFromContext(c).DeleteCampaignSchedule(campaignID)
}

// GetScheduledCampaigns returns all scheduled campaigns < time
func GetScheduledCampaigns(c context.Context, time time.Time) ([]entities.CampaignSchedule, error) {
	return GetFromContext(c).GetScheduledCampaigns(time)
}

// GetTotalSends returns total sends for specified campaign id
func GetTotalSends(c context.Context, campaignID, userID int64) (int64, error) {
	return GetFromContext(c).GetTotalSends(campaignID, userID)
}

// GetTotalDelivered returns  total delivered for specified campaign id
func GetTotalDelivered(c context.Context, campaignID, userID int64) (int64, error) {
	return GetFromContext(c).GetTotalDelivered(campaignID, userID)
}

// GetTotalBounces returns  total bounces for specified campaign id
func GetTotalBounces(c context.Context, campaignID, userID int64) (int64, error) {
	return GetFromContext(c).GetTotalBounces(campaignID, userID)
}

// GetTotalComplaints returns total complaints for specified campaign id
func GetTotalComplaints(c context.Context, campaignID, userID int64) (int64, error) {
	return GetFromContext(c).GetTotalComplaints(campaignID, userID)
}

// GetCampaignClicksStats returns a collection of clicks stats by given campaign id and user id
func GetCampaignClicksStats(c context.Context, id, userID int64) ([]entities.ClicksStats, error) {
	return GetFromContext(c).GetCampaignClicksStats(id, userID)
}

// GetSegments populates a pagination object with a collection of
// lists by the specified user id.
func GetSegments(c context.Context, userID int64, p *PaginationCursor) error {
	return GetFromContext(c).GetSegments(userID, p)
}

// GetSegmentsByIDs fetches lists by user id and the given ids
func GetSegmentsByIDs(c context.Context, userID int64, ids []int64) ([]entities.Segment, error) {
	return GetFromContext(c).GetSegmentsByIDs(userID, ids)
}

// GetSegment returns a Segment entity by the given id and user id.
func GetSegment(c context.Context, id, userID int64) (*entities.Segment, error) {
	return GetFromContext(c).GetSegment(id, userID)
}

// GetSegmentByName returns a SEgment entity by the given name and user id.
func GetSegmentByName(c context.Context, name string, userID int64) (*entities.Segment, error) {
	return GetFromContext(c).GetSegmentByName(name, userID)
}

// GetTotalSegments fetches the total count by user id.
func GetTotalSegments(c context.Context, userID int64) (int64, error) {
	return GetFromContext(c).GetTotalSegments(userID)
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
func GetSubscribers(c context.Context, userID int64, p *PaginationCursor, scopeMap map[string]string) error {
	return GetFromContext(c).GetSubscribers(userID, p, scopeMap)
}

// GetSubscribersBySegmentID populates a pagination object with a collection of
// subscribers by the specified user id and list id.
func GetSubscribersBySegmentID(c context.Context, segmentID, userID int64, p *PaginationCursor) error {
	return GetFromContext(c).GetSubscribersBySegmentID(segmentID, userID, p)
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

// DeactivateSubscriber blacklists a Subscriber entity by the given email.
func DeactivateSubscriber(c context.Context, userID int64, email string) error {
	return GetFromContext(c).DeactivateSubscriber(userID, email)
}

// DeleteSubscriber deletes a Subscriber entity by the given id.
func DeleteSubscriber(c context.Context, id, userID int64) error {
	return GetFromContext(c).DeleteSubscriber(id, userID)
}

// DeleteSubscriberByEmail deletes a Subscriber entity by the given email.
func DeleteSubscriberByEmail(c context.Context, email string, userID int64) error {
	return GetFromContext(c).DeleteSubscriberByEmail(email, userID)
}

// GetTotalSubscribers fetches the total count by user id.
func GetTotalSubscribers(c context.Context, userID int64) (int64, error) {
	return GetFromContext(c).GetTotalSubscribers(userID)
}

// GetTotalSubscribersBySegment fetches the total count by user and segment id.
func GetTotalSubscribersBySegment(c context.Context, segmentID, userID int64) (int64, error) {
	return GetFromContext(c).GetTotalSubscribersBySegment(segmentID, userID)
}

// SeekSubscribersByUserID returns subscribers for given user id
func SeekSubscribersByUserID(c context.Context, userID, nextID, limit int64) ([]entities.Subscriber, error) {
	return GetFromContext(c).SeekSubscribersByUserID(userID, nextID, limit)
}

// GetAPIKeys returns a list of APIKey entities.
func GetAPIKeys(c context.Context, userID int64) ([]*entities.APIKey, error) {
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

// GetToken returns the one-time token by the given token string.
// If the token is expired, we return an error indicating that a token is not found.
func GetToken(c context.Context, token string) (*entities.Token, error) {
	return GetFromContext(c).GetToken(token)
}

// CreateToken adds new token in the database.
func CreateToken(c context.Context, s *entities.Token) error {
	return GetFromContext(c).CreateToken(s)
}

// DeleteToken deletes the token by the given token.
func DeleteToken(c context.Context, token string) error {
	return GetFromContext(c).DeleteToken(token)
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

// CreateReport adds new report in the database.
func CreateReport(c context.Context, r *entities.Report) error {
	return GetFromContext(c).CreateReport(r)
}

// UpdateReport updates report in the database.
func UpdateReport(c context.Context, r *entities.Report) error {
	return GetFromContext(c).UpdateReport(r)
}

// GetReportByFilename returns report for provided user id and file name
func GetReportByFilename(c context.Context, filename string, userID int64) (*entities.Report, error) {
	return GetFromContext(c).GetReportByFilename(filename, userID)
}

// GetRunningReportForUser returns a report that is currently being generated for the specified user
func GetRunningReportForUser(c context.Context, userID int64) (*entities.Report, error) {
	return GetFromContext(c).GetRunningReportForUser(userID)
}

// GetNumberOfReportsForDate returns number of reports for date time.
func GetNumberOfReportsForDate(c context.Context, userID int64, time time.Time) (int64, error) {
	return GetFromContext(c).GetNumberOfReportsForDate(userID, time)
}

// GetTemplateByName returns a Template entity by the given name and user id.
func GetTemplateByName(c context.Context, name string, userID int64) (*entities.Template, error) {
	return GetFromContext(c).GetTemplateByName(name, userID)
}

// GetTemplate returns a Template entity by the given id and user id.
func GetTemplate(c context.Context, templateID, userID int64) (*entities.Template, error) {
	return GetFromContext(c).GetTemplate(templateID, userID)
}

// GetTemplates populates a pagination object with a collection of
// templates by the specified user id.
func GetTemplates(c context.Context, userID int64, p *PaginationCursor, scopeMap map[string]string) error {
	return GetFromContext(c).GetTemplates(userID, p, scopeMap)
}

// DeleteTemplate deletes the template with given template id and user id from db
func DeleteTemplate(c context.Context, templateID int64, userID int64) error {
	return GetFromContext(c).DeleteTemplate(templateID, userID)
}

// CreateSendLog creates a SendLogs entity.
func CreateSendLog(c context.Context, sendLogs *entities.SendLog) error {
	return GetFromContext(c).CreateSendLog(sendLogs)
}

// GetSendLogByUUID returns send log with specified uuid
func GetSendLogByUUID(c context.Context, id string) (*entities.SendLog, error) {
	return GetFromContext(c).GetSendLogByUUID(id)
}
