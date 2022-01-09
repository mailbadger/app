package storage

import (
	"time"

	"github.com/mailbadger/app/entities"
)

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
	LogFailedCampaign(c *entities.Campaign, description string) error

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

	CreateSendLog(l *entities.SendLog) error
	CountLogsByUUID(id string) (int64, error)
	CountLogsByStatus(status string) (int64, error)
	GetSendLogByUUID(id string) (*entities.SendLog, error)

	CreateBounce(b *entities.Bounce) error
	CreateComplaint(c *entities.Complaint) error
	CreateSend(s *entities.Send) error
	CreateClick(c *entities.Click) error
	CreateOpen(o *entities.Open) error
	CreateDelivery(d *entities.Delivery) error

	CreateReport(r *entities.Report) error
	UpdateReport(r *entities.Report) error
	GetReportByFilename(filename string, userID int64) (*entities.Report, error)
	GetRunningReportForUser(userID int64) (*entities.Report, error)
	GetNumberOfReportsForDate(userID int64, time time.Time) (int64, error)

	CreateTemplate(t *entities.Template) error
	UpdateTemplate(t *entities.Template) error
	GetTemplateByName(name string, userID int64) (*entities.Template, error)
	GetTemplate(templateID int64, userID int64) (*entities.Template, error)
	GetTemplates(userID int64, p *PaginationCursor, scopeMap map[string]string) error
	DeleteTemplate(templateID int64, userID int64) error
}
