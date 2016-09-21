package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func From(middleware ...gin.HandlerFunc) http.Handler {
	handler := gin.New()

	handler.Use(gin.Recovery())
	handler.Use(middleware...)

	return handler
}
