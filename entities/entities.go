package entities

import (
	"os"

	"bitbucket.org/liamstask/goose/lib/goose"
	"github.com/FilipNikolovski/news-maily/config"
	"github.com/FilipNikolovski/news-maily/utils"
	log "github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

//Global logger
var Logger = log.New()

var db *gorm.DB
var err error

//Opens database connection and runs the most recent migrations\
//If the db is not created it creates it and seeds it with a default user
func Setup() error {
	fresh_db := false
	if _, err = os.Stat(config.Config.Database); err != nil || config.Config.Database == ":memory:" {
		fresh_db = true
	}

	//Goose configuration
	migrateConfig := &goose.DBConf{
		MigrationsDir: config.Config.MigrationsDir,
		Env:           "production",
		Driver: goose.DBDriver{
			Name:    "sqlite3",
			OpenStr: config.Config.Database,
			Import:  "github.com/mattn/go-sqlite3",
			Dialect: &goose.Sqlite3Dialect{},
		},
	}

	//Get the most recent migration
	latest, err := goose.GetMostRecentDBVersion(migrateConfig.MigrationsDir)
	if err != nil {
		Logger.Println(err)
		return err
	}

	err = createDbConn()
	if err != nil {
		return err
	}

	//Run migrations
	err = goose.RunMigrationsOnDb(migrateConfig, migrateConfig.MigrationsDir, latest, db.DB())
	if err != nil {
		Logger.Println(err)
		return err
	}

	if fresh_db {
		// Hashing the password with the default cost of 10
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)
		if err != nil {
			Logger.Println(err)
			return err
		}

		key, err := utils.GenerateRandomString(32)
		if err != nil {
			Logger.Println(err)
			return err
		}

		//Create the default user
		admin := User{
			Username: "admin",
			Password: string(hashedPassword),
			ApiKey:   string(key),
		}

		err = db.Save(&admin).Error
		if err != nil {
			Logger.Println(err)
			return err
		}
	}
	return nil
}

//Creates a database connection to sqlite3
func createDbConn() error {
	db, err = gorm.Open("sqlite3", config.Config.Database)
	db.DB().SetMaxOpenConns(1)
	db.SetLogger(Logger)
	db.LogMode(false)

	if err != nil {
		Logger.Println(err)
		return err
	}

	return nil
}
