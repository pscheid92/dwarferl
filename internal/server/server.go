package server

import (
	"errors"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"github.com/pscheid92/dwarferl/internal"
	"github.com/pscheid92/dwarferl/internal/config"
	"net/http"
	"path/filepath"
	"regexp"
)

type Server struct {
	*gin.Engine

	// shared components
	Config       config.Configuration
	SessionStore sessions.Store

	// services
	Shortener internal.UrlShortenerService
	Users     internal.UsersService
}

func New(config config.Configuration, store sessions.Store, shortener internal.UrlShortenerService, users internal.UsersService) *Server {
	svr := &Server{
		Engine:       gin.New(),
		Config:       config,
		SessionStore: store,
		Shortener:    shortener,
		Users:        users,
	}

	goth.UseProviders(google.New(config.GoogleClientKey, config.GoogleSecret, config.GoogleCallbackURL))
	_ = svr.SetTrustedProxies(nil)
	svr.initHTMLRender()
	return svr
}

func (s *Server) initHTMLRender() {
	renderer := multitemplate.NewRenderer()

	// load the layout files
	layouts, err := filepath.Glob(s.Config.TemplatePath + "/layouts/*.gohtml")
	if err != nil {
		panic(err.Error())
	}

	// load the actual pages
	pages, err := filepath.Glob(s.Config.TemplatePath + "/*.gohtml")
	if err != nil {
		panic(err.Error())
	}

	// combine each actual page with layout files
	for _, page := range pages {
		name := filepath.Base(page)

		templates := make([]string, len(layouts)+1)
		templates[0] = page
		copy(templates[1:], layouts)

		renderer.AddFromFiles(name, templates...)
	}

	s.HTMLRender = renderer
}

func (s *Server) InitRoutes() {
	s.Use(sessions.Sessions("dwarferl_session", s.SessionStore))

	// public routes
	public := s.Group(s.Config.ForwardedPrefix)
	{
		public.Static("/assets", s.Config.AssetsPath)

		public.GET("/health", s.handleHealth())
		public.GET("/:short", s.handleRedirect())

		public.GET("/login", s.handleLoginPage())
		public.GET("/auth/:provider/callback", s.handleAuthCallback())
		public.GET("/auth/:provider", s.handleAuth())
		public.GET("/logout/:provider", s.handleLogout())
	}

	// private routes
	authorized := public.Group("")
	authorized.Use(s.authRequiredMiddleware())
	{
		authorized.GET("/", s.handleIndexPage())

		authorized.GET("/create", s.handleGetCreationPage())
		authorized.POST("/create", s.handlePostCreationPage())

		authorized.GET("/delete/:short", s.handleGetDeletionPage())
		authorized.POST("/delete/:short", s.handlePostDeletionPage())
	}
}

func (s *Server) handleHealth() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Status(http.StatusOK)
	}
}

func (s *Server) handleRedirect() gin.HandlerFunc {
	validShortRegex := regexp.MustCompilePOSIX(`^[A-Za-z0-9]+$`)

	return func(c *gin.Context) {
		short := c.Param("short")
		if short == "" || !validShortRegex.MatchString(short) {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		redirect, err := s.Shortener.ExpandShortURL(short)
		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.Header("Cache-Control", "private, max-age=90")
		c.Header("Referrer-Policy", "unsafe-url")
		c.Redirect(http.StatusMovedPermanently, redirect)
	}
}

func (s *Server) handleLoginPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")

		if userID != nil {
			c.Redirect(http.StatusFound, s.Config.ForwardedPrefix)
			return
		}

		data := gin.H{
			"userID":     userID,
			"linkPrefix": s.Config.ForwardedPrefix,
		}
		c.HTML(http.StatusOK, "login.gohtml", data)
	}
}

func (s *Server) handleAuthCallback() gin.HandlerFunc {
	return func(c *gin.Context) {
		externalUser, err := gothic.CompleteUserAuth(c.Writer, c.Request)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		user, err := s.Users.GetOrCreateByGoogle(externalUser.UserID, externalUser.Email)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		session := sessions.Default(c)
		session.Set("user_id", user.ID)
		if err := session.Save(); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.Redirect(http.StatusFound, s.Config.ForwardedPrefix)
	}
}

func (s *Server) handleLogout() gin.HandlerFunc {
	redirect := s.Config.ForwardedPrefix + "login"

	return func(c *gin.Context) {
		_ = gothic.Logout(c.Writer, c.Request)

		session := sessions.Default(c)
		session.Clear()
		if err := session.Save(); err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.Redirect(http.StatusTemporaryRedirect, redirect)
	}
}

func (s *Server) handleAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		q := c.Request.URL.Query()
		q.Add("provider", c.Param("provider"))
		c.Request.URL.RawQuery = q.Encode()

		user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
		if err != nil {
			gothic.BeginAuthHandler(c.Writer, c.Request)
			return
		}

		session := sessions.Default(c)
		session.Set("user_id", user.UserID)
		if err := session.Save(); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) handleIndexPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")

		list, err := s.Shortener.List(userID)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		data := gin.H{
			"redirects":  list,
			"userID":     userID,
			"linkPrefix": s.Config.ForwardedPrefix,
		}
		c.HTML(http.StatusOK, "index.gohtml", data)
	}
}

func (s *Server) handleGetCreationPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")

		data := gin.H{
			"userID":     userID,
			"linkPrefix": s.Config.ForwardedPrefix,
		}
		c.HTML(http.StatusOK, "create.gohtml", data)
	}
}

func (s *Server) handlePostCreationPage() gin.HandlerFunc {
	type request struct {
		Url string `form:"url"`
	}

	return func(c *gin.Context) {
		var req request
		if err := c.Bind(&req); err != nil {
			return
		}

		if req.Url == "" {
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("url is required"))
			return
		}

		userID := c.GetString("user_id")
		if _, err := s.Shortener.ShortenURL(req.Url, userID); err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.Redirect(http.StatusFound, s.Config.ForwardedPrefix)
	}
}

func (s *Server) handleGetDeletionPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		short := c.Param("short")
		userID := c.GetString("user_id")

		redirect, err := s.Shortener.GetRedirectByShort(short, userID)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		data := gin.H{
			"redirect":   redirect,
			"userID":     userID,
			"linkPrefix": s.Config.ForwardedPrefix,
		}
		c.HTML(http.StatusOK, "delete.gohtml", data)
	}
}

func (s *Server) handlePostDeletionPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		short := c.Param("short")
		userID := c.GetString("user_id")
		if err := s.Shortener.DeleteShortURL(short, userID); err != nil {
			_ = c.AbortWithError(http.StatusNotFound, err)
			return
		}
		c.Redirect(http.StatusFound, s.Config.ForwardedPrefix)
	}
}

func (s *Server) authRequiredMiddleware() gin.HandlerFunc {
	loginPage := s.Config.ForwardedPrefix + "login"

	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")
		if userID == nil {
			c.Redirect(http.StatusFound, loginPage)
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
