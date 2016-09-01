package entities

import (
	"log"
	"os"

	"bitbucket.org/liamstask/goose/lib/goose"
	"github.com/FilipNikolovski/news-maily/config"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

//Global logger
var Logger = log.New(os.Stdout, " ", log.Ldate|log.Ltime|log.Lshortfile)

var db *gorm.DB
var err error

//Opens database connection and runs the most recent migrations\
//If the db is not created it creates it and seeds it with a default user
func Setup() error {
	err = createDbConn()
	if err != nil {
		return err
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

	//Run migrations
	err = goose.RunMigrationsOnDb(migrateConfig, migrateConfig.MigrationsDir, latest, db.DB())
	if err != nil {
		Logger.Println(err)
		return err
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
