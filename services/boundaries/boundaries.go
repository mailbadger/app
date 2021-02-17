package boundaries

import (
	"fmt"

	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/storage"
)

// Service describes the boundaries interface for checking on resource limits.
type Service interface {
	CampaignsLimitExceeded(user *entities.User) (bool, error)
}

type service struct {
	store storage.Storage
}

// New returns a new boundaries service.
func New(store storage.Storage) Service {
	return &service{store}
}

func (svc *service) CampaignsLimitExceeded(user *entities.User) (bool, error) {
	limit := user.Boundaries.CampaignsLimit
	if limit > 0 {
		count, err := svc.store.GetTotalCampaigns(user.ID)
		if err != nil {
			return true, fmt.Errorf("boundaries: get total campaigns: %w", err)
		}

		return count >= limit, nil
	}

	return false, nil
}
