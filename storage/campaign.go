package storage

import (
	"github.com/FilipNikolovski/news-maily/entities"
	"github.com/FilipNikolovski/news-maily/utils/pagination"
)

// GetCampaigns fetches campaigns by user id, and populates the pagination obj
func (db *store) GetCampaigns(user_id int64, p *pagination.Pagination) {
	var campaigns []entities.Campaign
	var count uint64

	db.Offset(p.Offset).Limit(p.PerPage).Where("user_id = ?", user_id).Find(&campaigns).Count(&count)
	p.SetTotal(count)

	for _, t := range campaigns {
		p.Append(t)
	}
}

// GetCampaign returns the campaign by the given id and user id
func (db *store) GetCampaign(id int64, user_id int64) (*entities.Campaign, error) {
	var campaign = new(entities.Campaign)
	err := db.Where("user_id = ? and id = ?", user_id, id).Find(campaign).Error
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

// DeleteCampaign deletes an existing campaign in the database.
func (db *store) DeleteCampaign(id int64, user_id int64) error {
	return db.Where("user_id = ?", user_id).Delete(entities.Campaign{Id: id}).Error
}
