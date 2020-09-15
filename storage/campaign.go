package storage

import (
	"github.com/mailbadger/app/entities"
)

// GetCampaigns fetches campaigns by user id, and populates the pagination obj
func (db *store) GetCampaigns(userID int64, p *PaginationCursor) error {
	p.SetCollection(&[]entities.Campaign{})
	p.SetResource("campaigns")

	query := db.Table(p.Resource).
		Where("user_id = ?", userID).
		Order("created_at desc, id desc").
		Limit(p.PerPage)

	p.SetQuery(query)

	return db.Paginate(p, userID)
}

// GetTotalCampaigns fetches the total count by user id
func (db *store) GetTotalCampaigns(userID int64) (int64, error) {
	var count int64
	err := db.Model(entities.Campaign{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// GetCampaign returns the campaign by the given id and user id
func (db *store) GetCampaign(id, userID int64) (*entities.Campaign, error) {
	var campaign = new(entities.Campaign)
	err := db.Where("user_id = ? and id = ?", userID, id).Find(campaign).Error
	return campaign, err
}

// GetCampaignsByTemplateName returns a collection of campaigns by the given template id and user id
func (db *store) GetCampaignsByTemplateName(templateName string, userID int64) ([]entities.Campaign, error) {
	var campaigns []entities.Campaign
	err := db.Where("user_id = ? and template_name = ?", userID, templateName).Find(&campaigns).Error
	return campaigns, err
}

// GetCampaignByName returns the campaign by the given name and user id
func (db *store) GetCampaignByName(name string, userID int64) (*entities.Campaign, error) {
	var campaign = new(entities.Campaign)
	err := db.Where("user_id = ? and name = ?", userID, name).Find(campaign).Error
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
func (db *store) GetCampaignOpens(campaignID int64, p *PaginationCursor) error {
	p.SetCollection(&[]entities.Open{})
	p.SetResource("opens")

	query := db.Table(p.Resource).
		Where("campaign_id = ?", campaignID).
		Order("created_at desc, id desc").
		Limit(p.PerPage)

	p.SetQuery(query)

	return db.Paginate(p, campaignID)
}

// GetClicksStats fetches campaign total & unique clicks from the database.
func (db *store) GetClicksStats(campaignID int64) (*entities.ClicksStats, error) {
	clickStats := &entities.ClicksStats{}
	err := db.Table("clicks").Select("count(distinct(recipient))").Count(&clickStats.Unique).Select("count(recipient)").Count(&clickStats.Total).Where("campaign_id = ?", campaignID).Error
	return clickStats, err
}

// GetOpensStats fetches campaign total & unique opens from the database.
func (db *store) GetOpensStats(campaignID int64) (*entities.OpensStats, error) {
	opensStats := &entities.OpensStats{}
	err := db.Table("opens").Select("count(distinct(recipient))").Count(&opensStats.Unique).Select("count(recipient)").Count(&opensStats.Total).Where("campaign_id = ?", campaignID).Error
	return opensStats, err
}

// GetTotalSends returns total sends for campaign id from the database.
func (db *store) GetTotalSends(campaignID int64) (int64, error) {
	var totalSent int64
	err := db.Table("sends").Select("count(campaign_id)").Count(&totalSent).Where("campaign_id=?", campaignID).Error
	return totalSent, err
}

// GetTotalDelivered fetches campaign total deliveries  from the database.
func (db *store) GetTotalDelivered(campaignID int64) (int64, error) {
	var totalDelivered int64
	err := db.Table("deliveries").Select("count(distinct(recipient))").Count(&totalDelivered).Where("campaign_id = ?", campaignID).Error
	return totalDelivered, err
}

// GetTotalBounces fetches campaign total bounces  from the database.
func (db *store) GetTotalBounces(campaignID int64) (int64, error) {
	var totalBounces int64
	err := db.Table("bounces").Select("count(recipient)").Count(&totalBounces).Where("campaign_id = ?", campaignID).Error
	return totalBounces, err
}

// GetTotalComplaints fetches campaign total bounces  from the database.
func (db *store) GetTotalComplaints(campaignID int64) (int64, error) {
	var totalComplaints int64
	err := db.Table("complaints").Select("count(recipient)").Count(&totalComplaints).Where("campaign_id = ?", campaignID).Error
	return totalComplaints, err
}
