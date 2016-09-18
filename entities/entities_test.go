package entities

import (
	"encoding/base64"
	"errors"
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
	key, _ := base64.URLEncoding.DecodeString(user.ApiKey)

	assert.Nil(suite.Suite.T(), err)
	assert.Len(suite.Suite.T(), key, 32)
}

func (suite *EntitiesSuite) TestCreateTemplate() {
	err := CreateTemplate(&Template{
		Name:    "foo",
		UserId:  1,
		Content: "Foo bar",
	})

	suite.Suite.Equal(err, nil)
}

func (suite *EntitiesSuite) TestCreateTemplateNoName() {
	err := CreateTemplate(&Template{
		UserId:  1,
		Content: "Foo bar",
	})

	suite.Suite.Equal(err, errors.New("Name not specified"))
}

func (suite *EntitiesSuite) TestCreateTemplateNoContent() {
	err := CreateTemplate(&Template{
		Name:   "foo",
		UserId: 1,
	})

	suite.Suite.Equal(err, errors.New("Content not specified"))
}

func (suite *EntitiesSuite) TestGetTemplate() {
	CreateTemplate(&Template{
		Name:    "foo",
		UserId:  1,
		Content: "Foo bar",
	})

	t, err := GetTemplate(1, 1)
	suite.Suite.Equal(err, nil)
	suite.Suite.Equal(t.Name, "foo")
	suite.Suite.Equal(t.Content, "Foo bar")
}
