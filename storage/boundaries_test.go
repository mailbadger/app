package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoundaries(t *testing.T) {
	db := openTestDb()

	store := From(db)

	b, err := store.GetBoundariesByType("db_test")

	assert.Nil(t, err)
	assert.Equal(t, b.CampaignsLimit, int64(2))
	assert.True(t, b.SAMLEnabled)
	assert.True(t, b.ScheduleCampaignsEnabled)
}
