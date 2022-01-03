package config

import "github.com/kelseyhightower/envconfig"

type (
	Config struct {
		Database Database
		Redis    Redis
		Session  Session
		Server   Server
		Logging  Logging
		Consumer Consumer
		Mode     string `envconfig:"MB_APP_MODE"`
	}

	Database struct {
		Driver string `envconfig:"MB_APP_DATABASE_DRIVER" default:"sqlite3"`
		//Sqlite3
		Sqlite3Source string `envconfig:"MB_APP_SQLITE3_SOURCE" default:":memory:"`
		//Mysql
		MySQLUser     string `envconfig:"MB_APP_MYSQL_USER"`
		MySQLPass     string `envconfig:"MB_APP_MYSQL_PASS"`
		MySQLHost     string `envconfig:"MB_APP_MYSQL_HOST"`
		MySQLPort     string `envconfig:"MB_APP_MYSQL_PORT"`
		MySQLDatabase string `envconfig:"MB_APP_MYSQL_DATABASE"`
	}

	Redis struct {
		Host string `envconfig:"MB_APP_REDIS_HOST"`
		Port string `envconfig:"MB_APP_REDIS_PORT"`
		Pass string `envconfig:"MB_APP_REDIS_PASS"`
	}

	Session struct {
		Secure     bool   `envconfig:"MB_APP_SECURE_COOKIE"`
		AuthKey    string `envconfig:"MB_APP_SESSION_AUTH_KEY"`
		EncryptKey string `envconfig:"MB_APP_SESSION_ENCRYPT_KEY"`
	}

	Server struct {
		Port                string `envconfig:"MB_APP_PORT" default:":8082"`
		Cert                string `envconfig:"MB_APP_TLS_CERT"`
		Key                 string `envconfig:"MB_APP_TLS_KEY"`
		AppDir              string `envconfig:"MB_APP_DIR"`
		AppURL              string `envconfig:"MB_APP_URL"`
		UnsubscribeSecret   string `envconfig:"MB_APP_UNSUBSCRIBE_SECRET"`
		SystemEmailSource   string `envconfig:"MB_APP_SYSTEM_EMAIL_SOURCE"`
		EnableSignup        bool   `envconfig:"MB_APP_ENABLE_SIGNUP"`
		VerifyEmailOnSignup bool   `envconfig:"MB_APP_VERIFY_EMAIL_ON_SIGNUP"`
		RecaptchaSecret     string `envconfig:"MB_APP_RECAPTCHA_SECRET"`
	}

	Logging struct {
		Level  string `envconfig:"MB_APP_LOG_LEVEL" default:"info"`
		Pretty bool   `envconfig:"MB_APP_LOG_PRETTY"`
	}

	Consumer struct {
		Timeout         int32 `envconfig:"MB_APP_CONSUMER_TIMEOUT" default:"300"`
		WaitTimeout     int32 `envconfig:"MB_APP_CONSUMER_WAIT_TIMEOUT" default:"10"`
		MaxInFlightMsgs int32 `envconfig:"MB_APP_CONSUMER_MAX_INFLIGHT_MSGS" default:"10"`
	}
)

// FromEnv returns the Config object from the environment.
func FromEnv() (Config, error) {
	c := Config{}
	err := envconfig.Process("", &c)
	return c, err
}
