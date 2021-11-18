package storage

import (
	"testing"
	"time"

	"github.com/mailbadger/app/entities"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestTokens(t *testing.T) {
	db := openTestDb()
	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			logrus.Error(err)
		}
		sqlDB.Close()
	}()
	store := From(db)

	_, err := store.GetToken("abc")
	assert.NotNil(t, err)

	now := time.Now()
	token := &entities.Token{
		UserID:    1,
		Token:     "abc",
		Type:      entities.UnsubscribeTokenType,
		ExpiresAt: now.AddDate(0, 0, 7),
	}

	err = store.CreateToken(token)
	assert.Nil(t, err)

	token, err = store.GetToken("abc")
	assert.Nil(t, err)
	assert.Equal(t, "abc", token.Token)
	assert.Equal(t, entities.UnsubscribeTokenType, token.Type)

	err = store.DeleteToken("abc")
	assert.Nil(t, err)

	token, err = store.GetToken("abc")
	assert.NotNil(t, err)
	assert.Nil(t, token)

	// Create expired token
	token = &entities.Token{
		UserID:    1,
		Token:     "abc",
		Type:      entities.UnsubscribeTokenType,
		ExpiresAt: now.AddDate(0, 0, -1),
	}

	err = store.CreateToken(token)
	assert.Nil(t, err)

	token, err = store.GetToken("abc")
	assert.NotNil(t, err)
	assert.Nil(t, token)

	// test delete all tokens for user
	token = &entities.Token{
		UserID:    1,
		Token:     "delete-all-tokens",
		Type:      entities.UnsubscribeTokenType,
		ExpiresAt: now.AddDate(0, 0, 4),
	}

	err = store.CreateToken(token)
	assert.Nil(t, err)

	err = store.DeleteAllTokensForUser(1)
	assert.Nil(t, err)

	token, err = store.GetToken("delete-all-tokens")
	assert.NotNil(t, err)
	assert.Empty(t, token)
}
