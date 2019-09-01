package storage

import (
	"time"

	"github.com/news-maily/app/entities"
	"github.com/news-maily/app/utils/pagination"
	"github.com/sirupsen/logrus"
)

// GetCampaigns fetches campaigns by user id, and populates the pagination obj
func (db *store) GetCampaigns(userID int64, p *pagination.Cursor) {
	var campaigns []entities.Campaign

	var reverse bool
	var prevID, nextID int64

	query := db.Where("user_id = ?", userID).Limit(p.PerPage).Order("created_at desc, id desc")

	if p.EndingBefore != 0 {
		c, err := db.GetCampaign(p.EndingBefore, userID)
		if err != nil {
			logrus.WithFields(logrus.Fields{"ending_before": p.EndingBefore, "user_id": userID}).WithError(err).
				Error("Unable to find campaign for pagination with ending before id.")
			return
		}

		query.Where(`(created_at > ? OR (created_at = ? AND id > ?)) AND created_at < ?`,
			c.CreatedAt.Format(time.RFC3339Nano),
			c.CreatedAt.Format(time.RFC3339Nano),
			c.ID,
			time.Now().Format(time.RFC3339Nano),
		).
			Order("created_at, id", true).Find(&campaigns)

		// populate prev and next
		if len(campaigns) > 0 {
			nextID = campaigns[0].ID
			last, err := db.getLastCampaign(userID)
			if err != nil {
				logrus.WithFields(logrus.Fields{"user_id": userID}).WithError(err).
					Error("Unable to find the last campaign.")
				return
			}

			if last.ID != campaigns[len(campaigns)-1].ID {
				prevID = campaigns[len(campaigns)-1].ID
			}
		}

		reverse = true
	} else if p.StartingAfter != 0 {
		c, err := db.GetCampaign(p.StartingAfter, userID)
		if err != nil {
			logrus.WithFields(logrus.Fields{"starting_after": p.StartingAfter, "user_id": userID}).WithError(err).
				Error("Unable to find campaign for pagination with starting after id.")
			return
		}
		query.Where(`(created_at < ? OR (created_at = ? AND id < ?)) AND created_at < ?`,
			c.CreatedAt.Format(time.RFC3339Nano),
			c.CreatedAt.Format(time.RFC3339Nano),
			c.ID,
			time.Now().Format(time.RFC3339Nano),
		).Find(&campaigns)

		// populate prev and next
		if len(campaigns) > 0 {
			prevID = campaigns[0].ID
			first, err := db.getFirstCampaign(userID)
			if err != nil {
				logrus.WithFields(logrus.Fields{"user_id": userID}).WithError(err).
					Error("Unable to find the first campaign.")
				return
			}

			if first.ID != campaigns[len(campaigns)-1].ID {
				nextID = campaigns[len(campaigns)-1].ID
			}
		}
	} else {
		query.Find(&campaigns)
		if len(campaigns) > 0 {
			nextID = campaigns[len(campaigns)-1].ID
		}
	}

	if reverse {
		for i := len(campaigns) - 1; i >= 0; i-- {
			p.Append(campaigns[i])
		}
	} else {
		for _, s := range campaigns {
			p.Append(s)
		}
	}

	p.PopulateLinks(prevID, nextID)
}

func (db *store) getFirstCampaign(userID int64) (*entities.Campaign, error) {
	var c = new(entities.Campaign)
	err := db.Where("user_id = ?", userID).Order("created_at, id").Limit(1).Find(c).Error
	return c, err
}

func (db *store) getLastCampaign(userID int64) (*entities.Campaign, error) {
	var c = new(entities.Campaign)
	err := db.Where("user_id = ?", userID).Order("created_at desc, id desc").Limit(1).Find(c).Error
	return c, err
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
	return db.Where("user_id = ?", userID).Delete(entities.Campaign{ID: id}).Error
}
