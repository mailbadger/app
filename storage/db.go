package storage

import (
	"database/sql"
	"fmt"
	"os"
	"time"
	
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rakyll/statik/fs"
	migrate "github.com/rubenv/sql-migrate"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	
	"github.com/mailbadger/app/mode"
	_ "github.com/mailbadger/app/statik"
	
	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/utils"
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
func openDbConn(driver, config string) (db *gorm.DB) {
	var err error
	
	gormConfig := &gorm.Config{}
	if mode.IsDebug() {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}
	
	switch driver {
	case "mysql":
		db, err = gorm.Open(mysql.Open(config), gormConfig)
	
	case "sqlite3":
		db, err = gorm.Open(sqlserver.Open(config), gormConfig)
	}
	
	if err != nil {
		log.WithError(err).Fatalln("db connection failed")
	}
	
	sqlDB, err := db.DB()
	if err != nil {
		log.WithError(err).Fatalln("sql db connection failed")
	}
	
	sqlDB.SetMaxIdleConns(0)
	
	if err := pingDb(sqlDB); err != nil {
		log.WithError(err).Fatalln("database ping attempts failed")
	}
	
	if err := setupDb(driver, config, db); err != nil {
		log.WithError(err).Fatalln("migrations failed")
	}
	
	return db
}

// setupDb runs the necessary migrations and creates a new user if the database
// hasn't been setup yet
func setupDb(driver, config string, db *gorm.DB) error {
	log.Info("Running migrations..")
	
	migrationFS, err := fs.NewWithNamespace("migrations")
	if err != nil {
		return fmt.Errorf("create migrations file system: %w", err)
	}
	
	var m = &migrate.HttpFileSystemMigrationSource{
		FileSystem: migrationFS,
	}
	
	sqlDB, err := db.DB()
	if err != nil {
		log.WithError(err).Fatalln("sql db connection failed")
	}
	
	_, err = migrate.Exec(sqlDB, driver, m, migrate.Up)
	if err != nil {
		return err
	}
	
	// If the database didn't exist, initialize it with an admin user
	if isFresh(driver, config, db) {
		err = initDb(config, db)
		if err != nil {
			return err
		}
	}
	log.Info("DB is up to date..")
	
	return nil
}

// isFresh is checking if the database is empty
func isFresh(driver, config string, db *gorm.DB) bool {
	switch driver {
	case "sqlite3":
		if _, err := os.Stat(config); err != nil || config == ":memory:" {
			return true
		}
	case "mysql":
		err := db.First(&entities.User{}).Error
		if err != nil {
			return true
		}
	}
	
	return false
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
	
	// Create the default user
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

// MakeConfigFromEnv creates a DSN string from env variables based on the driver name.
// List of drivers: 'sqlite3', 'mysql'.
func MakeConfigFromEnv(driver string) string {
	switch driver {
	case "sqlite3":
		return os.Getenv("SQLITE3_FILE")
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true",
			os.Getenv("MYSQL_USER"),
			os.Getenv("MYSQL_PASS"),
			os.Getenv("MYSQL_HOST"),
			os.Getenv("MYSQL_PORT"),
			os.Getenv("MYSQL_DATABASE"),
		)
	default:
		return ""
	}
}
