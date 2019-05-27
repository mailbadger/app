package storage

import (
	"database/sql"
	"os"
	"time"

	"github.com/google/uuid"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/news-maily/app/entities"
	"github.com/news-maily/app/storage/migrations"
	"github.com/news-maily/app/utils"
	migrate "github.com/rubenv/sql-migrate"
	log "github.com/sirupsen/logrus"
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

// From creates a new store object.
func From(db *gorm.DB) Storage {
	return &store{db}
}

// openDbConn creates a database connection using the driver and config string
func openDbConn(driver, config string) *gorm.DB {
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

	fresh := false
	switch driver {
	case "sqlite3":
		if _, err := os.Stat(config); err != nil || config == ":memory:" {
			fresh = true
		}
	case "mysql":
		err := db.First(&entities.User{}).Error
		if err != nil {
			fresh = true
		}
	}

	if err := setupDb(driver, config, fresh, db); err != nil {
		log.Errorln(err)
		log.Fatalln("migrations failed")
	}

	// db.LogMode(utils.IsDebugMode())

	return db
}

// setupDb runs the necessary migrations and creates a new user if the database
// hasn't been setup yet
func setupDb(driver, config string, fresh bool, db *gorm.DB) error {
	log.Info("Running migrations..")

	var m = &migrate.AssetMigrationSource{
		Asset:    migrations.Asset,
		AssetDir: migrations.AssetDir,
		Dir:      driver,
	}
	_, err := migrate.Exec(db.DB(), driver, m, migrate.Up)
	if err != nil {
		log.Errorln(err)
		return err
	}

	// If the database didn't exist, initialize it with an admin user
	if fresh {
		err = initDb(config, db)
		if err != nil {
			log.Errorln(err)
			return err
		}
	}
	log.Info("DB is up to date..")

	return nil
}

// initDb seeds the database with the admin user, if the database has not been
// initialized before
func initDb(config string, db *gorm.DB) error {
	log.Info("Generating new credentials...")

	// Hashing the password with the default cost of 10
	secret, err := utils.GenerateRandomString(12)
	if err != nil {
		log.Println(err)
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return err
	}

	uuid, err := uuid.NewRandom()
	if err != nil {
		log.Errorf("unable to generate random uuid: %s", err.Error())
		return err
	}

	//Create the default user
	admin := entities.User{
		Username: "admin",
		UUID:     uuid.String(),
		Password: sql.NullString{
			String: string(hashedPassword),
			Valid:  true,
		},
		Active:   true,
		Verified: true,
		Source:   "mailbadger.io",
	}

	err = db.Save(&admin).Error
	if err != nil {
		log.Errorln(err)
		return err
	}

	log.WithFields(log.Fields{"user": "admin", "password": secret}).Info("Admin user credentials")

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
