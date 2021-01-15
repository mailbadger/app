package limits

import (
	"context"

	"github.com/mailbadger/app/storage"
)

type service struct {
	store storage.Storage
}

func (svc *service) Exceeded(ctx context.Context) (bool, error) {
	return false, nil
}

func (svc *service) SubscriptionExpired(ctx context.Context) bool {
	return false
}
