package boundaries

import (
	"bytes"
	"fmt"
	"io"

	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/storage"
)

// Service describes the boundaries interface for checking on resource limits.
type Service interface {
	CampaignsLimitExceeded(user *entities.User) (bool, error)
	SubscribersLimitExceeded(user *entities.User) (bool, int64, error)
	CSVLineCounter(r io.Reader) (int, error)
}

type service struct {
	store storage.Storage
}

// New returns a new boundaries service.
func New(store storage.Storage) Service {
	return &service{store}
}

func (s *service) CampaignsLimitExceeded(user *entities.User) (bool, error) {
	limit := user.Boundaries.CampaignsLimit
	if limit > 0 {
		count, err := s.store.GetMonthlyTotalCampaigns(user.ID)
		if err != nil {
			return true, fmt.Errorf("boundaries: get total campaigns: %w", err)
		}

		return count >= limit, nil
	}

	return false, nil
}

func (s *service) SubscribersLimitExceeded(user *entities.User) (bool, int64, error) {
	limit := user.Boundaries.SubscribersLimit
	if limit > 0 {
		count, err := s.store.GetTotalSubscribers(user.ID)
		if err != nil {
			return true, 0, fmt.Errorf("boundaries: get total subscribers: %w", err)
		}
		return count >= limit, count, err
	}
	return false, 0, nil
}

func (s *service) CSVLineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}
