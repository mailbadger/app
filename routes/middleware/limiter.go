package middleware

import (
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/didip/tollbooth_gin"
	"github.com/gin-gonic/gin"
)

func Limiter() gin.HandlerFunc {
	lmt := tollbooth.NewLimiter(10, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
	lmt.SetMessage(`{"message": "You have reached the maximum request limit."}`)
	lmt.SetMessageContentType("application/json; charset=utf-8")
	return tollbooth_gin.LimitHandler(lmt)
}
