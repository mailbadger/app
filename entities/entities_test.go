package entities

import (
	"testing"

	"github.com/FilipNikolovski/news-maily/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type EntitiesSuite struct {
	suite.Suite

	User *User
}

func (suite *EntitiesSuite) SetupSuite() {
	Logger.Info("Setting up..")
	config.Config.Database = ":memory:"
	config.Config.MigrationsDir = "../db/migrations/"
	err := Setup()
	if err != nil {
		Logger.Error("Failed creating db: %v", err)
	}
}

func TestEntitiesSuite(t *testing.T) {
	suite.Run(t, new(EntitiesSuite))
}

func (suite *EntitiesSuite) TestGetUser() {
	user, err := GetUser(1)
	suite.Suite.Equal(err, nil)
	suite.Suite.Equal(user.Username, "admin")
}

func (suite *EntitiesSuite) TestUpdateUser() {
	user, err := GetUser(1)
	user.Username = "foo"
	err = UpdateUser(&user)
	suite.Suite.Equal(err, nil)

	user, err = GetUser(1)
	suite.Suite.Equal(user.Username, "foo")
}

func (suite *EntitiesSuite) TestApiKeyLength() {
	user, err := GetUser(1)
	assert.Nil(suite.Suite.T(), err)
	assert.Len(suite.Suite.T(), user.ApiKey, 16)
}
