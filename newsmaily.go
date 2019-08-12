package main

/*
news-maily - Open-Source Newsletter Mail System

The MIT License (MIT)

Copyright (c) 2015 Filip Nikolovski

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is furnished
to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included
in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/news-maily/app/routes"
	"github.com/news-maily/app/utils"
	"github.com/sirupsen/logrus"
)

func init() {
	lvl, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		lvl = logrus.InfoLevel
	}

	logrus.SetLevel(lvl)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)

	if !utils.IsDebugMode() {
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

	if os.Getenv("ENVIRONMENT") == "prod" && os.Getenv("SKIP_TLS") == "" {
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
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
			Certificates: []tls.Certificate{cer},
		}
	}

	srv := &http.Server{
		Addr:         ":" + addr,
		Handler:      handler,
		TLSConfig:    cfg,
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
			logrus.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	logrus.Infoln("Starting HTTP server on port", srv.Addr)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		logrus.Printf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
}
