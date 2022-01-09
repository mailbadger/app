package storage

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/mailbadger/app/config"
	"github.com/mailbadger/app/entities"
	_ "github.com/mailbadger/app/statik"
	"github.com/mailbadger/app/utils"
	"github.com/rakyll/statik/fs"
	migrate "github.com/rubenv/sql-migrate"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// store implements the Storage interface
type store struct {
	*gorm.DB
}

// New creates a database connection and returns a new DB
func New(conf config.Config) *gorm.DB {
	dsn := makeDsn(conf)
	return openDbConn(conf.Storage.DB.Driver, dsn)
}

// From creates a new store object.
func From(db *gorm.DB) Storage {
	return &store{db}
}

// openDbConn creates a database connection using the driver and source string
func openDbConn(driver, dsn string) *gorm.DB {
	var dialect gorm.Dialector
	if driver == "mysql" {
		dialect = mysql.Open(dsn)
	} else {
		dialect = sqlite.Open(dsn)
	}

	conf := &gorm.Config{}

	db, err := gorm.Open(dialect, conf)
	if err != nil {
		log.WithError(err).Fatalln("db open connection failed")
	}

	conn, err := db.DB()
	if err != nil {
		log.WithError(err).Fatalln("db get connection failed")
	}

	if driver == "mysql" {
		conn.SetMaxIdleConns(0)
	}

	if err := pingDb(conn); err != nil {
		log.WithError(err).Fatalln("database ping attempts failed")
	}

	fresh := false
	switch driver {
	case "sqlite3":
		if _, err := os.Stat(dsn); err != nil || dsn == ":memory:" {
			fresh = true
		}
	case "mysql":
		err := db.First(&entities.User{}).Error
		if err != nil {
			fresh = true
		}
	}

	if err := setupDb(driver, dsn, fresh, db); err != nil {
		log.WithError(err).Fatalln("database setup failed")
	}

	return db
}

// pingDb ensures that the database is reachable before running migrations
func pingDb(db *sql.DB) (err error) {
	for i := 0; i < 20; i++ {
		err = db.Ping()
		if err == nil {
			return
		}

		log.Infof("database ping failed. retry in 1s")
		time.Sleep(time.Second)
	}
	return
}

// makeDsn creates a DSN string from the db config based on the driver name.
// List of drivers: 'sqlite3', 'mysql'.
func makeDsn(conf config.Config) string {
	switch conf.Storage.DB.Driver {
	case "sqlite3":
		return conf.Storage.DB.Sqlite3Source
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true",
			conf.Storage.DB.MySQLUser,
			conf.Storage.DB.MySQLPass,
			conf.Storage.DB.MySQLHost,
			conf.Storage.DB.MySQLPort,
			conf.Storage.DB.MySQLDatabase,
		)
	default:
		return ""
	}
}

// setupDb runs the necessary migrations and creates a new user if the database
// hasn't been setup yet
func setupDb(driver, config string, fresh bool, db *gorm.DB) error {
	log.Info("Running migrations..")

	migrationFS, err := fs.NewWithNamespace("migrations")
	if err != nil {
		return fmt.Errorf("create migrations file system: %w", err)
	}

	var m = &migrate.HttpFileSystemMigrationSource{
		FileSystem: migrationFS,
	}
	conn, err := db.DB()
	if err != nil {
		return fmt.Errorf("create migrations db conn: %w", err)
	}
	_, err = migrate.Exec(conn, driver, m, migrate.Up)
	if err != nil {
		return err
	}

	// If the database didn't exist, initialize it with an admin user
	if fresh {
		err = initDb(config, db)
		if err != nil {
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
		return fmt.Errorf("init db: gen rand string: %w", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("init db: hash password: %w", err)
	}

	uuid := uuid.New()

	//Create the default user
	nolimit := &entities.Boundaries{}
	err = db.Where("type = ?", entities.BoundaryTypeNoLimit).First(nolimit).Error
	if err != nil {
		return fmt.Errorf("init db: fetch nolimit boundaries: %w", err)
	}

	adminRole := entities.Role{}
	err = db.Where("name = ?", entities.AdminRole).First(&adminRole).Error
	if err != nil {
		return fmt.Errorf("init db: fetch admin role: %w", err)
	}

	admin := entities.User{
		Username: "admin",
		UUID:     uuid.String(),
		Password: sql.NullString{
			String: string(hashedPassword),
			Valid:  true,
		},
		Active:     true,
		Verified:   true,
		Boundaries: nolimit,
		Roles:      []entities.Role{adminRole},
		Source:     "mailbadger.io",
	}

	err = db.Save(&admin).Error
	if err != nil {
		return fmt.Errorf("init db: save user: %w", err)
	}

	log.WithFields(log.Fields{"user": "admin", "password": secret}).Info("Admin user credentials..make sure to change that password!")

	return nil
}

// openTestDb creates a database connection for testing purposes
func openTestDb() *gorm.DB {
	var (
		driver = "sqlite3"
		config = ":memory:"
	)

	return openDbConn(driver, config)
}
