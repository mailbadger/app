package routes

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/didip/tollbooth_gin"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/contrib/ginrus"
	"github.com/gin-gonic/gin"
	"github.com/open-policy-agent/opa/ast"
	"github.com/sirupsen/logrus"
	"github.com/unrolled/secure"

	"github.com/mailbadger/app/actions"
	"github.com/mailbadger/app/mode"
	"github.com/mailbadger/app/opa"
	"github.com/mailbadger/app/routes/middleware"
	"github.com/mailbadger/app/s3"
	"github.com/mailbadger/app/storage"
	"github.com/mailbadger/app/templates"
)

// New creates a new HTTP handler with the specified middleware.
func New() http.Handler {
	lvl, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		lvl = logrus.InfoLevel
	}

	log := logrus.New()
	log.SetLevel(lvl)
	log.SetOutput(os.Stdout)
	if mode.IsProd() {
		log.SetFormatter(&logrus.JSONFormatter{})
	}

	s3Client, err := s3.NewS3Client(
		os.Getenv("AWS_S3_ACCESS_KEY"),
		os.Getenv("AWS_S3_SECRET_KEY"),
		os.Getenv("AWS_S3_REGION"),
	)
	if err != nil {
		panic(err)
	}

	store := cookie.NewStore(
		[]byte(os.Getenv("SESSION_AUTH_KEY")),
		[]byte(os.Getenv("SESSION_ENCRYPT_KEY")),
	)
	secureCookie, _ := strconv.ParseBool(os.Getenv("SECURE_COOKIE"))
	store.Options(sessions.Options{
		Secure:   secureCookie,
		HttpOnly: true,
	})

	driver := os.Getenv("DATABASE_DRIVER")
	config := storage.MakeConfigFromEnv(driver)

	s := storage.New(driver, config)

	handler := gin.New()

	handler.Use(gin.Recovery())
	handler.Use(ginrus.Ginrus(log, time.RFC3339, true))
	handler.Use(sessions.Sessions("mbsess", store))
	handler.Use(middleware.Storage(s))
	handler.Use(middleware.Producer())
	handler.Use(middleware.SetUser())
	handler.Use(middleware.RequestID())
	handler.Use(middleware.SetLoggerEntry())
	handler.Use(middleware.S3Client(s3Client))

	// Security headers
	secureMiddleware := secure.New(secure.Options{
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		SSLRedirect:           true,
		SSLProxyHeaders:       map[string]string{"X-Forwarded-Proto": "https"},
		STSSeconds:            31536000,
		STSIncludeSubdomains:  true,
		STSPreload:            true,
		ContentSecurityPolicy: "default-src 'self';style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; font-src 'self' https://fonts.gstatic.com; script-src 'self' 'unsafe-inline'",

		IsDevelopment: !mode.IsProd(),
	})
	secureFunc := func() gin.HandlerFunc {
		return func(c *gin.Context) {
			err := secureMiddleware.Process(c.Writer, c.Request)

			// If there was an error, do not continue.
			if err != nil {
				c.Abort()
				return
			}

			// Avoid header rewrite if response is a redirection.
			if status := c.Writer.Status(); status > 300 && status < 399 {
				c.Abort()
			}
		}
	}()

	handler.Use(secureFunc)

	// Web app
	appDir := os.Getenv("APP_DIR")
	if appDir == "" {
		logrus.Panic("app directory not set")
	}

	err = templates.Init(handler)
	if err != nil {
		logrus.Panic(err)
	}

	handler.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Not found.",
			})
			return
		}

		if strings.HasPrefix(c.Request.URL.Path, "/unsubscribe.html") {
			email := c.Query("email")
			t := c.Query("t")
			uuid := c.Query("uuid")
			failed := c.Query("failed")

			c.HTML(http.StatusOK, "unsubscribe.html", gin.H{
				"email":  email,
				"t":      t,
				"uuid":   uuid,
				"failed": failed,
			})
			return
		}

		if strings.HasPrefix(c.Request.URL.Path, "/unsubscribe-success.html") {
			c.HTML(http.StatusOK, "unsubscribe-success.html", nil)
			return
		}

		c.File(appDir + "/index.html")
	})

	// Assets
	handler.Static("/static", appDir+"/static")

	// rate limiter
	lmt := tollbooth.NewLimiter(10, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
	lmt.SetMessage(`{"message": "You have reached the maximum request limit."}`)
	lmt.SetMessageContentType("application/json; charset=utf-8")

	SetGuestRoutes(
		handler,
		middleware.NoCache(),
		tollbooth_gin.LimitHandler(lmt),
	)

	// Compile the OPA module. The keys are used as identifiers in error messages.
	opacompiler, err := opa.NewCompiler()
	if err != nil {
		panic(err)
	}

	SetAuthorizedRoutes(
		handler,
		opacompiler,
		middleware.NoCache(),
		middleware.CSRF(),
		tollbooth_gin.LimitHandler(lmt),
	)

	return handler
}

// SetGuestRoutes sets the guest routes to the gin engine handler along with
// a number of middleware that we set.
func SetGuestRoutes(handler *gin.Engine, middleware ...gin.HandlerFunc) {
	guest := handler.Group("/api")
	guest.Use(middleware...)

	guest.GET("/auth/github/callback", actions.GithubCallback)
	guest.GET("/auth/github", actions.GetGithubAuth)
	guest.GET("/auth/google/callback", actions.GoogleCallback)
	guest.GET("/auth/google", actions.GetGoogleAuth)
	guest.GET("/auth/facebook/callback", actions.FacebookCallback)
	guest.GET("/auth/facebook", actions.GetFacebookAuth)
	guest.POST("/authenticate", actions.PostAuthenticate)
	guest.POST("/forgot-password", actions.PostForgotPassword)
	guest.PUT("/forgot-password/:token", actions.PutForgotPassword)
	guest.PUT("/verify-email/:token", actions.PutVerifyEmail)
	guest.POST("/signup", actions.PostSignup)
	guest.POST("/hooks/:uuid", actions.HandleHook)
	guest.POST("/unsubscribe", actions.PostUnsubscribe)
}

// SetAuthorizedRoutes sets the authorized routes to the gin engine handler along with
// the Authorized middleware which performs the checks for authorized user as well as
// other optional middlewares that we set.
func SetAuthorizedRoutes(handler *gin.Engine, opacompiler *ast.Compiler, middlewares ...gin.HandlerFunc) {
	authorized := handler.Group("/api")
	authorized.Use(middleware.Authorized(opacompiler))
	authorized.Use(middlewares...)

	authorized.POST("/logout", actions.PostLogout)
	{
		users := authorized.Group("/users")
		{
			users.GET("/me", actions.GetMe)
			users.POST("/password", actions.ChangePassword)
		}

		templates := authorized.Group("/templates")
		{
			templates.GET("", middleware.PaginateWithCursor(), actions.GetTemplates)
			templates.GET("/:id", actions.GetTemplate)
			templates.POST("", actions.PostTemplate)
			templates.PUT("/:id", actions.PutTemplate)
			templates.DELETE("/:id", actions.DeleteTemplate)
		}

		campaigns := authorized.Group("/campaigns")
		{
			campaigns.GET("", middleware.PaginateWithCursor(), actions.GetCampaigns)
			campaigns.GET("/:id", actions.GetCampaign)
			campaigns.POST("", actions.PostCampaign)
			campaigns.PUT("/:id", actions.PutCampaign)
			campaigns.DELETE("/:id", actions.DeleteCampaign)
			campaigns.POST("/:id/start", actions.StartCampaign)
			campaigns.GET("/:id/opens", middleware.PaginateWithCursor(), actions.GetCampaignOpens)
			campaigns.GET("/:id/stats", actions.GetCampaignStats)
			campaigns.GET("/:id/clicks", actions.GetCampaignClicksStats)
			campaigns.GET("/:id/complaints", middleware.PaginateWithCursor(), actions.GetCampaignComplaints)
			campaigns.GET("/:id/bounces", middleware.PaginateWithCursor(), actions.GetCampaignBounces)
			campaigns.PATCH("/:id/schedule", actions.PatchCampaignSchedule)
			campaigns.DELETE("/:id/schedule", actions.DeleteCampaignSchedule)
		}

		segments := authorized.Group("/segments")
		{
			segments.GET("", middleware.PaginateWithCursor(), actions.GetSegments)
			segments.GET("/:id", actions.GetSegment)
			segments.POST("", actions.PostSegment)
			segments.PUT("/:id", actions.PutSegment)
			segments.DELETE("/:id", actions.DeleteSegment)
			segments.PUT("/:id/subscribers", actions.PutSegmentSubscribers)
			segments.GET("/:id/subscribers", middleware.PaginateWithCursor(), actions.GetSegmentsubscribers)
			segments.POST("/:id/subscribers/detach", actions.DetachSegmentSubscribers)
			segments.DELETE("/:id/subscribers/:sub_id", actions.DetachSubscriber)
		}

		subscribers := authorized.Group("/subscribers")
		{
			subscribers.GET("", middleware.PaginateWithCursor(), actions.GetSubscribers)
			subscribers.GET("/:id", actions.GetSubscriber)
			subscribers.GET("/export/download", actions.DownloadSubscribersReport)
			subscribers.POST("", actions.PostSubscriber)
			subscribers.PUT("/:id", actions.PutSubscriber)
			subscribers.DELETE("/:id", actions.DeleteSubscriber)
			subscribers.POST("/import", actions.ImportSubscribers)
			subscribers.POST("/bulk-remove", actions.BulkRemoveSubscribers)
			subscribers.POST("/export", actions.ExportSubscribers)
		}

		ses := authorized.Group(("/ses"))
		{
			ses.GET("/keys", actions.GetSESKeys)
			ses.POST("/keys", actions.PostSESKeys)
			ses.DELETE("/keys", actions.DeleteSESKeys)
			ses.GET("/quota", actions.GetSESQuota)
		}

		s3 := authorized.Group("/s3")
		{
			s3.POST("/sign", actions.GetSignedURL)
		}
	}
}
