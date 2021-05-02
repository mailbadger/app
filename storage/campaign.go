package storage

import (
	"os"

	"github.com/jinzhu/gorm"

	"github.com/mailbadger/app/entities"
)

// GetCampaigns fetches campaigns by user id, and populates the pagination obj
func (db *store) GetCampaigns(userID int64, p *PaginationCursor, scopeMap map[string]string) error {
	p.SetCollection(&[]entities.Campaign{})
	p.SetResource("campaigns")

	// scopes
	p.AddScope(NotDeleted)
	for k, v := range scopeMap {
		if k == "name" {
			p.AddScope(NameLike(v))
		}
	}

	query := db.Table(p.Resource).Preload("BaseTemplate").
		Where("user_id = ?", userID).
		Order("created_at desc, id desc").
		Limit(p.PerPage)

	p.SetQuery(query)

	return db.Paginate(p, userID)
}

// GetMonthlyTotalCampaigns fetches the total count by user id in the current month
func (db *store) GetMonthlyTotalCampaigns(userID int64) (int64, error) {
	var count int64
	err := db.Model(entities.Campaign{}).
		Scopes(currentMonthScope()).
		Where("user_id = ?", userID).
		Count(&count).Error
	return count, err
}

func currentMonthScope() func(db *gorm.DB) *gorm.DB {
	driver := os.Getenv("DATABASE_DRIVER")
	switch driver {
	case "mysql":
		return currentMonthMySQLDialect
	default:
		return currentMonthSqlite3Dialect
	}
}

func currentMonthMySQLDialect(db *gorm.DB) *gorm.DB {
	return db.Where("YEAR(created_at) = YEAR(NOW()) AND MONTH(created_at) = MONTH(NOW())")
}

func currentMonthSqlite3Dialect(db *gorm.DB) *gorm.DB {
	return db.Where("strftime('%Y', created_at) = strftime('%Y', date('now')) AND strftime('%m', created_at) = strftime('%m',date('now'))")
}

// GetCampaign returns the campaign by the given id and user id
func (db *store) GetCampaign(id, userID int64) (*entities.Campaign, error) {
	var campaign = new(entities.Campaign)
	err := db.Where("user_id = ? and id = ?", userID, id).Preload("BaseTemplate").Preload("Schedule").Find(&campaign).Error
	return campaign, err
}

// GetCampaignByName returns the campaign by the given name and user id
func (db *store) GetCampaignByName(name string, userID int64) (*entities.Campaign, error) {
	var campaign = new(entities.Campaign)
	err := db.Preload("BaseTemplate").Where("user_id = ? and name = ?", userID, name).Find(campaign).Error
	return campaign, err
}

// CreateCampaign creates a new campaign in the database.
func (db *store) CreateCampaign(c *entities.Campaign) error {
	return db.Create(c).Error
}

// UpdateCampaign edits an existing campaign in the database.
func (db *store) UpdateCampaign(c *entities.Campaign) error {
	return db.Where("id = ? and user_id = ?", c.ID, c.UserID).Save(c).Error
}

// DeleteCampaign deletes an existing campaign from the database.
func (db *store) DeleteCampaign(id, userID int64) error {
	return db.Where("user_id = ?", userID).Delete(entities.Campaign{Model: entities.Model{ID: id}}).Error
}

// GetCampaignOpens fetches campaign opens by campaign id, and populates the pagination obj
func (db *store) GetCampaignOpens(campaignID, userID int64, p *PaginationCursor) error {
	p.SetCollection(&[]entities.Open{})
	p.SetResource("opens")

	query := db.Table(p.Resource).
		Where("campaign_id = ? and user_id=?", campaignID, userID).
		Order("created_at desc, id desc").
		Limit(p.PerPage)

	p.SetQuery(query)

	return db.Paginate(p, userID)
}

// GetClicksStats fetches campaign total & unique clicks from the database.
func (db *store) GetClicksStats(campaignID, userID int64) (*entities.ClicksStats, error) {
	clickStats := &entities.ClicksStats{}
	err := db.Table("clicks").
		Where("campaign_id = ? and user_id= ?", campaignID, userID).
		Select("count(distinct(recipient))").Count(&clickStats.UniqueClicks).
		Select("count(recipient)").Count(&clickStats.TotalClicks).Error
	return clickStats, err
}

// GetOpensStats fetches campaign total & unique opens from the database.
func (db *store) GetOpensStats(campaignID, userID int64) (*entities.OpensStats, error) {
	opensStats := &entities.OpensStats{}
	err := db.Table("opens").
		Where("campaign_id = ? and user_id= ?", campaignID, userID).
		Select("count(distinct(recipient))").Count(&opensStats.Unique).
		Select("count(recipient)").Count(&opensStats.Total).Error
	return opensStats, err
}

// GetTotalSends returns total sends for campaign id from the database.
func (db *store) GetTotalSends(campaignID, userID int64) (int64, error) {
	var totalSent int64
	err := db.Table("sends").
		Where("campaign_id = ? and user_id= ?", campaignID, userID).Count(&totalSent).Error
	return totalSent, err
}

// GetTotalDelivered fetches campaign total deliveries  from the database.
func (db *store) GetTotalDelivered(campaignID, userID int64) (int64, error) {
	var totalDelivered int64
	err := db.Table("deliveries").
		Where("campaign_id = ? and user_id= ?", campaignID, userID).Count(&totalDelivered).Error
	return totalDelivered, err
}

// GetTotalBounces fetches campaign total bounces  from the database.
func (db *store) GetTotalBounces(campaignID, userID int64) (int64, error) {
	var totalBounces int64
	err := db.Table("bounces").
		Where("campaign_id = ? and user_id= ?", campaignID, userID).Count(&totalBounces).Error
	return totalBounces, err
}

// GetTotalComplaints fetches campaign total bounces  from the database.
func (db *store) GetTotalComplaints(campaignID, userID int64) (int64, error) {
	var totalComplaints int64
	err := db.Table("complaints").
		Where("campaign_id = ? and user_id= ?", campaignID, userID).Count(&totalComplaints).Error
	return totalComplaints, err
}

// GetCampaignComplaints fetches campaign complaints by campaign id, and populates the pagination obj
func (db *store) GetCampaignComplaints(campaignID, userID int64, p *PaginationCursor) error {
	p.SetCollection(&[]entities.Complaint{})
	p.SetResource("complaints")

	query := db.Table(p.Resource).
		Where("campaign_id = ? and user_id = ?", campaignID, userID).
		Order("created_at desc, id desc").
		Limit(p.PerPage)

	p.SetQuery(query)

	return db.Paginate(p, userID)
}

// GetCampaignBounces fetches campaign bounces by campaign id, and populates the pagination obj
func (db *store) GetCampaignBounces(campaignID, userID int64, p *PaginationCursor) error {
	p.SetCollection(&[]entities.Bounce{})
	p.SetResource("bounces")

	query := db.Table(p.Resource).
		Where("campaign_id = ? and user_id= ?", campaignID, userID).
		Order("created_at desc, id desc").
		Limit(p.PerPage)

	p.SetQuery(query)

	return db.Paginate(p, userID)
}

// DeleteAllCampaignsForUser deletes all campaigns for user
func (db *store) DeleteAllCampaignsForUser(userID int64) error {
	return db.Where("user_id = ?", userID).Delete(&entities.Campaign{}).Error
}

// NameLike applies a scope for campaigns by the given name.
// The wildcard is applied on the end of the name search.
func NameLike(name string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("name LIKE ?", name+"%")
	}
}

// NotDeleted scopes a resource by deletion column.
func NotDeleted(db *gorm.DB) *gorm.DB {
	return db.Where("deleted_at IS NULL")
}
