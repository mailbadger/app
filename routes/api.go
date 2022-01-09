package routes

import (
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/mailbadger/app/actions"
	"github.com/mailbadger/app/config"
	"github.com/mailbadger/app/emails"
	"github.com/mailbadger/app/routes/middleware"
	"github.com/mailbadger/app/services/boundaries"
	"github.com/mailbadger/app/services/reports"
	"github.com/mailbadger/app/services/subscribers"
	templatesvc "github.com/mailbadger/app/services/templates"
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
	sqsPublisher sqs.PublisherAPI
	s3Client     s3iface.S3API
	emailSender  emails.Sender
	templatesvc  templatesvc.Service
	boundarysvc  boundaries.Service
	subscrsvc    subscribers.Service
	reportsvc    reports.Service

	campaignerQueueURL sqs.CampaignerQueueURL
	appDir             string
	appURL             string

	filesBucket string

	enableSignup           bool
	verifyEmail            bool
	recaptchaSecret        string
	unsubscribeTokenSecret string
	systemEmail            string
	social                 config.Social
}

func From(
	sess session.Session,
	store storage.Storage,
	opaCompiler *ast.Compiler,
	sqsPublisher sqs.PublisherAPI,
	s3Client s3iface.S3API,
	emailSender emails.Sender,
	templatesvc templatesvc.Service,
	boundarysvc boundaries.Service,
	subscrsvc subscribers.Service,
	reportsvc reports.Service,
	campaignerQueueURL sqs.CampaignerQueueURL,
	conf config.Config,
) API {
	return New(
		sess,
		store,
		opaCompiler,
		sqsPublisher,
		s3Client,
		emailSender,
		templatesvc,
		boundarysvc,
		subscrsvc,
		reportsvc,
		campaignerQueueURL,
		conf.Server.AppDir,
		conf.Server.AppURL,
		conf.Storage.S3.FilesBucket,
		conf.Server.EnableSignup,
		conf.Server.VerifyEmailOnSignup,
		conf.Server.RecaptchaSecret,
		conf.Server.UnsubscribeSecret,
		conf.Server.SystemEmailSource,
		conf.Social,
	)
}
func New(
	sess session.Session,
	store storage.Storage,
	opaCompiler *ast.Compiler,
	sqsPublisher sqs.PublisherAPI,
	s3Client s3iface.S3API,
	emailSender emails.Sender,
	templatesvc templatesvc.Service,
	boundarysvc boundaries.Service,
	subscrsvc subscribers.Service,
	reportsvc reports.Service,
	campaignerQueueURL sqs.CampaignerQueueURL,
	appDir string,
	appURL string,
	filesBucket string,
	enableSignup bool,
	verifyEmail bool,
	recaptchaSecret string,
	unsubscribeTokenSecret string,
	systemEmail string,
	social config.Social,
) API {
	return API{
		sess:                   sess,
		store:                  store,
		opaCompiler:            opaCompiler,
		sqsPublisher:           sqsPublisher,
		s3Client:               s3Client,
		emailSender:            emailSender,
		templatesvc:            templatesvc,
		boundarysvc:            boundarysvc,
		subscrsvc:              subscrsvc,
		reportsvc:              reportsvc,
		campaignerQueueURL:     campaignerQueueURL,
		appDir:                 appDir,
		appURL:                 appURL,
		filesBucket:            filesBucket,
		enableSignup:           enableSignup,
		verifyEmail:            verifyEmail,
		recaptchaSecret:        recaptchaSecret,
		unsubscribeTokenSecret: unsubscribeTokenSecret,
		systemEmail:            systemEmail,
		social:                 social,
	}
}

func (api API) Handler() http.Handler {
	handler := gin.New()
	handler.Use(gin.Recovery())
	handler.Use(middleware.Secure())
	handler.Use(middleware.RequestID())
	handler.Use(middleware.Logger())
	handler.Use(sessions.Sessions("mbsess", api.sess.CookieStore))

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

	api.SetGuestRoutes(
		handler,
		middleware.NoCache(),
		middleware.Limiter(),
	)

	api.SetAuthorizedRoutes(
		handler,
		middleware.NoCache(),
		middleware.Limiter(),
		middleware.CSRF(api.sess.AuthKey, api.sess.Secure),
	)

	return handler
}

// SetGuestRoutes sets the guest routes to the gin engine handler along with
// a number of middleware that we set.
func (api API) SetGuestRoutes(handler *gin.Engine, middleware ...gin.HandlerFunc) {
	guest := handler.Group("/api")
	guest.Use(middleware...)

	guest.GET("/auth/github/callback",
		actions.GithubCallback(
			api.store,
			api.sess,
			api.social.Github.ClientID,
			api.social.Github.ClientSecret,
			api.appURL,
		),
	)
	guest.GET("/auth/github", actions.GetGithubAuth(api.social.Github.ClientID))

	guest.GET("/auth/google/callback",
		actions.GoogleCallback(
			api.store,
			api.sess,
			api.social.Google.ClientID,
			api.social.Google.ClientSecret,
			api.appURL,
		),
	)
	guest.GET("/auth/google",
		actions.GetGoogleAuth(
			api.social.Google.ClientID,
			api.social.Google.ClientSecret,
			api.appURL,
		),
	)

	guest.GET("/auth/facebook",
		actions.GetFacebookAuth(
			api.social.Facebook.ClientID,
			api.appURL,
		),
	)
	guest.GET("/auth/facebook/callback",
		actions.FacebookCallback(
			api.store,
			api.sess,
			api.social.Facebook.ClientID,
			api.social.Facebook.ClientSecret,
			api.appURL,
		))

	guest.POST("/authenticate", actions.PostAuthenticate(api.store, api.sess))
	guest.POST("/forgot-password",
		actions.PostForgotPassword(
			api.store,
			api.emailSender,
			api.systemEmail,
			api.appURL,
		),
	)
	guest.PUT("/forgot-password/:token", actions.PutForgotPassword(api.store))
	guest.PUT("/verify-email/:token", actions.PutVerifyEmail(api.store))
	guest.POST("/signup",
		actions.PostSignup(
			api.store,
			api.sess,
			api.emailSender,
			api.enableSignup,
			api.verifyEmail,
			api.recaptchaSecret,
			api.systemEmail,
			api.appURL,
		),
	)
	guest.POST("/hooks/:uuid", actions.HandleHook(api.store))
	guest.POST("/unsubscribe",
		actions.PostUnsubscribe(
			api.store,
			api.unsubscribeTokenSecret,
			api.appURL,
		),
	)
}

// SetAuthorizedRoutes sets the authorized routes to the gin engine handler along with
// the Authorized middleware which performs the checks for authorized user as well as
// other optional middlewares that we set.
func (api API) SetAuthorizedRoutes(handler *gin.Engine, middlewares ...gin.HandlerFunc) {
	authorized := handler.Group("/api")
	authorized.Use(middleware.Authorized(api.sess, api.store, api.opaCompiler))
	authorized.Use(middlewares...)

	authorized.POST("/logout", actions.PostLogout(api.sess))
	{
		users := authorized.Group("/users")
		{
			users.GET("/me", actions.GetMe)
			users.POST("/password", actions.ChangePassword(api.store))
		}

		templates := authorized.Group("/templates")
		{
			templates.GET("", middleware.PaginateWithCursor(), actions.GetTemplates(api.templatesvc))
			templates.GET("/:id", actions.GetTemplate(api.templatesvc))
			templates.POST("", actions.PostTemplate(api.templatesvc, api.store))
			templates.PUT("/:id", actions.PutTemplate(api.templatesvc, api.store))
			templates.DELETE("/:id", actions.DeleteTemplate(api.templatesvc))
		}

		campaigns := authorized.Group("/campaigns")
		{
			campaigns.GET("", middleware.PaginateWithCursor(), actions.GetCampaigns(api.store))
			campaigns.GET("/:id", actions.GetCampaign(api.store))
			campaigns.POST("", actions.PostCampaign(api.boundarysvc, api.store))
			campaigns.PUT("/:id", actions.PutCampaign(api.store))
			campaigns.DELETE("/:id", actions.DeleteCampaign(api.store))
			campaigns.POST("/:id/start", actions.StartCampaign(api.store, api.sqsPublisher, api.campaignerQueueURL))
			campaigns.GET("/:id/opens", middleware.PaginateWithCursor(), actions.GetCampaignOpens(api.store))
			campaigns.GET("/:id/stats", actions.GetCampaignStats(api.store))
			campaigns.GET("/:id/clicks", actions.GetCampaignClicksStats(api.store))
			campaigns.GET("/:id/complaints", middleware.PaginateWithCursor(), actions.GetCampaignComplaints(api.store))
			campaigns.GET("/:id/bounces", middleware.PaginateWithCursor(), actions.GetCampaignBounces(api.store))
			campaigns.PATCH("/:id/schedule", actions.PatchCampaignSchedule(api.store))
			campaigns.DELETE("/:id/schedule", actions.DeleteCampaignSchedule(api.store))
		}

		segments := authorized.Group("/segments")
		{
			segments.GET("", middleware.PaginateWithCursor(), actions.GetSegments(api.store))
			segments.GET("/:id", actions.GetSegment(api.store))
			segments.POST("", actions.PostSegment(api.store))
			segments.PUT("/:id", actions.PutSegment(api.store))
			segments.DELETE("/:id", actions.DeleteSegment(api.store))
			segments.PUT("/:id/subscribers", actions.PutSegmentSubscribers(api.store))
			segments.GET("/:id/subscribers", middleware.PaginateWithCursor(), actions.GetSegmentsubscribers(api.store))
			segments.POST("/:id/subscribers/detach", actions.DetachSegmentSubscribers(api.store))
			segments.DELETE("/:id/subscribers/:sub_id", actions.DetachSubscriber(api.store))
		}

		subscribers := authorized.Group("/subscribers")
		{
			subscribers.GET("", middleware.PaginateWithCursor(), actions.GetSubscribers(api.store))
			subscribers.GET("/:id", actions.GetSubscriber(api.store))
			subscribers.GET("/export/download", actions.DownloadSubscribersReport(api.store, api.s3Client, api.filesBucket))
			subscribers.POST("", actions.PostSubscriber(api.boundarysvc, api.store))
			subscribers.PUT("/:id", actions.PutSubscriber(api.store))
			subscribers.DELETE("/:id", actions.DeleteSubscriber(api.store))
			subscribers.POST("/import", actions.ImportSubscribers(
				api.subscrsvc,
				api.boundarysvc,
				api.store,
				api.s3Client,
				api.filesBucket,
			))
			subscribers.POST("/bulk-remove", actions.BulkRemoveSubscribers(api.subscrsvc, api.s3Client, api.filesBucket))
			subscribers.POST("/export", actions.ExportSubscribers(api.reportsvc, api.filesBucket))
		}

		ses := authorized.Group(("/ses"))
		{
			ses.GET("/keys", actions.GetSESKeys(api.store))
			ses.POST("/keys", actions.PostSESKeys(api.store, api.appURL))
			ses.DELETE("/keys", actions.DeleteSESKeys(api.store))
			ses.GET("/quota", actions.GetSESQuota(api.store))
		}

		s3 := authorized.Group("/s3")
		{
			s3.POST("/sign", actions.GetSignedURL(api.s3Client, api.filesBucket))
		}
	}
}
