package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mailbadger/app/mode"
)

//nolint
func initMode(m string) {
	mode.SetMode(m)
	if mode.IsProd() {
		gin.SetMode(gin.ReleaseMode)
	}
}
