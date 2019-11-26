package storage

import (
	"testing"
	"time"

	"github.com/news-maily/app/entities"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestTokens(t *testing.T) {
	db := openTestDb()
	defer func() {
		err := db.Close()
		if err != nil {
			logrus.Error(err)
		}
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

	//Create expired token
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
}
