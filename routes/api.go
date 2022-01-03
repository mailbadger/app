package routes

import (
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/mailbadger/app/actions"
	"github.com/mailbadger/app/config"
	"github.com/mailbadger/app/routes/middleware"
	"github.com/mailbadger/app/session"
	"github.com/mailbadger/app/sqs"
	"github.com/mailbadger/app/storage"
	"github.com/mailbadger/app/templates"
	"github.com/open-policy-agent/opa/ast"
	"github.com/sirupsen/logrus"
)

type API struct {
	sess         session.Session
	store        storage.Storage
	opaCompiler  *ast.Compiler
	sqsPublisher sqs.Publisher
	s3Client     *s3.S3
	appDir       string
}

func From(
	sess session.Session,
	store storage.Storage,
	opaCompiler *ast.Compiler,
	sqsPublisher sqs.Publisher,
	s3Client *s3.S3,
	conf config.Config,
) API {
	return New(
		sess,
		store,
		opaCompiler,
		sqsPublisher,
		s3Client,
		conf.Server.AppDir,
	)
}
func New(
	sess session.Session,
	store storage.Storage,
	opaCompiler *ast.Compiler,
	sqsPublisher sqs.Publisher,
	s3Client *s3.S3,
	appDir string,
) API {
	return API{
		sess:         sess,
		store:        store,
		opaCompiler:  opaCompiler,
		sqsPublisher: sqsPublisher,
		s3Client:     s3Client,
		appDir:       appDir,
	}
}

func (api API) Handler() http.Handler {
	handler := gin.New()
	handler.Use(gin.Recovery())
	handler.Use(middleware.Secure())
	handler.Use(middleware.RequestID())
	handler.Use(middleware.Logger())
	handler.Use(sessions.Sessions("mbsess", api.sess.CookieStore))
	handler.Use(middleware.Storage(api.store))
	handler.Use(middleware.S3Client(api.s3Client))
	handler.Use(middleware.SQSPublisher(api.sqsPublisher))

	err := templates.Init(handler)
	if err != nil {
		logrus.WithError(err).Panic("api: unable to init templates")
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

		c.File(api.appDir + "/index.html")
	})

	handler.Static("/static", api.appDir+"/static")

	SetGuestRoutes(
		handler,
		middleware.NoCache(),
		middleware.Limiter(),
	)

	SetAuthorizedRoutes(
		handler,
		api.sess,
		api.opaCompiler,
		middleware.NoCache(),
		middleware.Limiter(),
		middleware.CSRF(),
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
func SetAuthorizedRoutes(handler *gin.Engine, sess session.Session, opacompiler *ast.Compiler, middlewares ...gin.HandlerFunc) {
	authorized := handler.Group("/api")
	authorized.Use(middleware.Authorized(sess, opacompiler))
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
