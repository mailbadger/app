package session

import (
	"errors"
	"fmt"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/mailbadger/app/config"
	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/utils"
)

type Store interface {
	GetSession(id string) (*entities.Session, error)
	CreateSession(sess *entities.Session) error
}

type Session struct {
	Secure      bool
	AuthKey     string
	EncryptKey  string
	CookieStore cookie.Store

	store Store
}

const (
	sessKey      = "sess_id"
	sessDuration = 72 * time.Hour
)

var (
	ErrNotFound         = errors.New("session not found")
	ErrInvalidValueType = errors.New("session has invalid value type")
)

func From(store Store, conf config.Config) Session {
	return New(
		store,
		conf.Session.AuthKey,
		conf.Session.EncryptKey,
		conf.Session.Secure,
	)
}

func New(store Store, authKey, encryptKey string, secure bool) Session {
	cookieStore := cookie.NewStore(
		[]byte(authKey),
		[]byte(encryptKey),
	)
	cookieStore.Options(sessions.Options{
		Secure:   secure,
		HttpOnly: true,
	})

	return Session{
		store:       store,
		CookieStore: cookieStore,
		Secure:      secure,
		AuthKey:     authKey,
		EncryptKey:  encryptKey,
	}
}

func (sess Session) GetSession(c *gin.Context) (*entities.Session, error) {
	defaultsess := sessions.Default(c)
	v := defaultsess.Get(sessKey)
	if v == nil {
		return nil, ErrNotFound
	}
	sessID, ok := v.(string)
	if !ok {
		return nil, ErrInvalidValueType
	}
	s, err := sess.store.GetSession(sessID)
	return s, err
}

func (sess Session) CreateSession(c *gin.Context, userID int64) error {
	sessID, err := utils.GenerateRandomString(32)
	if err != nil {
		return fmt.Errorf("session: gen session id: %w", err)
	}

	err = sess.store.CreateSession(&entities.Session{
		UserID:    userID,
		SessionID: sessID,
	})
	if err != nil {
		return fmt.Errorf("session: create session: %w", err)
	}

	session := sessions.Default(c)
	exp := time.Now().Add(sessDuration).Unix() - time.Now().Unix()

	session.Options(sessions.Options{
		HttpOnly: true,
		MaxAge:   int(exp),
		Secure:   sess.Secure,
		Path:     "/api",
	})
	session.Set(sessKey, sessID)

	err = session.Save()
	if err != nil {
		return fmt.Errorf("session: save: %w", err)
	}
	return nil
}
