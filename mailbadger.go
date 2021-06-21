package main

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/mailbadger/app/mode"
	"github.com/mailbadger/app/routes"
)

func init() {
	mode.SetModeFromEnv()
	lvl, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		lvl = logrus.InfoLevel
	}

	logrus.SetLevel(lvl)
	logrus.SetOutput(os.Stdout)
	if mode.IsProd() {
		logrus.SetFormatter(&logrus.JSONFormatter{})
		gin.SetMode(gin.ReleaseMode)
	}
}

func main() {
	handler := routes.New()

	var cfg *tls.Config
	var addr = os.Getenv("PORT")
	if addr == "" {
		addr = "8080"
	}
	skipTLS, _ := strconv.ParseBool(os.Getenv("SKIP_TLS"))

	if !skipTLS {
		// TLS config
		cer, err := tls.LoadX509KeyPair(os.Getenv("CERT_FILE"), os.Getenv("KEY_FILE"))
		if err != nil {
			logrus.WithError(err).Error("unable to load x509 key pair")
			return
		}
		cfg = &tls.Config{
			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305, // Go 1.8 only
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,   // Go 1.8 only
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			},
			Certificates: []tls.Certificate{cer},
		}
	}

	srv := &http.Server{
		Addr:         ":" + addr,
		Handler:      handler,
		TLSConfig:    cfg,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)

		// interrupt signal sent from terminal
		signal.Notify(sigint, os.Interrupt)
		// sigterm signal sent from kubernetes
		signal.Notify(sigint, syscall.SIGTERM)

		<-sigint

		// We received an interrupt signal, shut down.
		if err := srv.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			logrus.Infof("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	logrus.Infoln("Starting HTTP server on port", srv.Addr)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		logrus.Infof("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
}
