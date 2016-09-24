package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// New creates a new HTTP handler with the specified middleware.
func New(middleware ...gin.HandlerFunc) http.Handler {
	handler := gin.New()

	handler.Use(gin.Recovery())
	handler.Use(middleware...)

	return handler
}
