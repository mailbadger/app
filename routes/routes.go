package routes

import (
	"net/http"
	"time"

	"github.com/FilipNikolovski/news-maily/routes/middleware"
	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/contrib/ginrus"
	"github.com/gin-gonic/gin"
)

// New creates a new HTTP handler with the specified middleware.
func New() http.Handler {
	handler := gin.New()

	handler.Use(gin.Recovery())
	handler.Use(ginrus.Ginrus(logrus.StandardLogger(), time.RFC3339, true))
	handler.Use(middleware.Storage())
	handler.Use(middleware.SetUser())

	return handler
}
