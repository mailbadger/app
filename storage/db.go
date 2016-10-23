package storage

import (
	"os"
	"time"

	"bitbucket.org/liamstask/goose/lib/goose"
	"github.com/FilipNikolovski/news-maily/entities"
	"github.com/FilipNikolovski/news-maily/utils"
	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

// store implements the Storage interface
type store struct {
	*gorm.DB
}

// New creates a database connection and returns a new Storage
func New(driver, config string) Storage {
	return From(openDbConn(driver, config))
}

func From(db *gorm.DB) Storage {
	return &store{db}
}

// openDbConn creates a database connection using the driver and config string
func openDbConn(driver, config string) *gorm.DB {
	fresh := false
	if _, err := os.Stat(config); err != nil || config == ":memory:" {
		fresh = true
	}

	db, err := gorm.Open(driver, config)
	if err != nil {
		log.Errorln(err)
		log.Fatalln("db connection failed!")
	}

	if driver == "mysql" {
		db.DB().SetMaxIdleConns(0)
	}

	if err := pingDb(db); err != nil {
		log.Errorln(err)
		log.Fatalln("database ping attempts failed")
	}

	if err := setupDb(driver, config, fresh, db); err != nil {
		log.Errorln(err)
		log.Fatalln("migrations failed")
	}

	return db
}

// setupDb runs the necessary migrations and creates a new user if the database
// hasn't been setup yet
func setupDb(driver, config string, fresh bool, db *gorm.DB) error {
	//Goose configuration
	migrateConfig := &goose.DBConf{
		MigrationsDir: "./migrations/sqlite3",
		Driver: goose.DBDriver{
			Name:    driver,
			OpenStr: config,
			Dialect: &goose.Sqlite3Dialect{},
		},
	}
	//Get the most recent migration
	latest, err := goose.GetMostRecentDBVersion(migrateConfig.MigrationsDir)
	if err != nil {
		log.Println(err)
	}

	//Run migrations
	err = goose.RunMigrationsOnDb(migrateConfig, migrateConfig.MigrationsDir, latest, db.DB())
	if err != nil {
		log.Errorln(err)
	}

	// If the database didn't exist, initialize it with an admin user
	if fresh {
		err = initDb(config, db)
		if err != nil {
			log.Errorln(err)
			return err
		}
	}
	return nil
}

// initDb seeds the database with the admin user, if the database has not been
// initialized before
func initDb(config string, db *gorm.DB) error {
	// Hashing the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return err
	}

	key, err := utils.GenerateRandomString(32)
	if err != nil {
		log.Errorln(err)
		return err
	}

	//Create the default user
	admin := entities.User{
		Username: "admin",
		Password: string(hashedPassword),
		ApiKey:   string(key),
	}

	err = db.Save(&admin).Error
	if err != nil {
		log.Errorln(err)
		return err
	}

	return nil
}

// pingDb ensures that the database is reachable before running migrations
func pingDb(db *gorm.DB) (err error) {
	for i := 0; i < 20; i++ {
		err = db.DB().Ping()
		if err == nil {
			return
		}

		log.Infof("database ping failed. retry in 1s")
		time.Sleep(time.Second)
	}
	return
}

// openTestDb creates a database connection for testing purposes
func openTestDb() *gorm.DB {
	var (
		driver = "sqlite3"
		config = ":memory:"
	)

	if os.Getenv("DATABASE_DRIVER") != "" && os.Getenv("DATABASE_CONFIG") != "" {
		driver = os.Getenv("DATABASE_DRIVER")
		config = os.Getenv("DATABASE_CONFIG")
	}

	return openDbConn(driver, config)
}
