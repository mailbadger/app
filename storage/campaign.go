package storage

import (
	"github.com/news-maily/api/entities"
	"github.com/news-maily/api/utils/pagination"
)

// GetCampaigns fetches campaigns by user id, and populates the pagination obj
func (db *store) GetCampaigns(userID int64, p *pagination.Pagination) {
	var campaigns []entities.Campaign
	var count uint64

	db.Offset(p.Offset).Limit(p.PerPage).Where("user_id = ?", userID).Find(&campaigns).Count(&count)
	p.SetTotal(count)

	for _, t := range campaigns {
		p.Append(t)
	}
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
	return db.Where("id = ? and user_id = ?", c.Id, c.UserId).Save(c).Error
}

// DeleteCampaign deletes an existing campaign from the database.
func (db *store) DeleteCampaign(id, userID int64) error {
	return db.Where("user_id = ?", userID).Delete(entities.Campaign{Id: id}).Error
}
