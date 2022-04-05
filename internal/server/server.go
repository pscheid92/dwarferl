package server

import (
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"github.com/pscheid92/dwarferl/internal"
	"github.com/pscheid92/dwarferl/internal/config"
	"net/http"
	"path"
)

type Server struct {
	*gin.Engine

	// shared components
	Config       config.Configuration
	SessionStore sessions.Store

	// services
	Shortener internal.UrlShortenerService
}

func New(config config.Configuration, store sessions.Store, shortener internal.UrlShortenerService) *Server {
	svr := &Server{
		Engine:       gin.New(),
		Config:       config,
		SessionStore: store,
		Shortener:    shortener,
	}

	goth.UseProviders(google.New(config.GoogleClientKey, config.GoogleSecret, config.GoogleCallbackURL))

	svr.LoadHTMLGlob(path.Join(config.TemplatePath, "*.gohtml"))
	return svr
}

func (s *Server) InitRoutes() {
	s.Use(sessions.Sessions("dwarferl_session", s.SessionStore))

	// public routes
	public := s.Group(s.Config.ForwardedPrefix)
	{
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
	return func(c *gin.Context) {
		short := c.Param("short")

		redirect, err := s.Shortener.ExpandShortURL(short)
		if err != nil {
			c.AbortWithStatus(404)
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
		c.HTML(http.StatusOK, "login.gohtml", userID)
	}
}

func (s *Server) handleAuthCallback() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		session := sessions.Default(c)
		session.Set("user_id", user.UserID)
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
		list, err := s.Shortener.List("00000000-0000-0000-0000-000000000000")
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		data := gin.H{
			"redirects":  list,
			"linkPrefix": s.Config.ForwardedPrefix,
		}

		c.HTML(http.StatusOK, "index.gohtml", data)
	}
}

func (s *Server) handleGetCreationPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "create.gohtml", nil)
	}
}

func (s *Server) handlePostCreationPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			Url string `form:"url"`
		}

		err := c.Bind(&request)
		if err != nil {
			return
		}

		if request.Url == "" {
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("url is required"))
			return
		}

		_, err = s.Shortener.ShortenURL(request.Url)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.Redirect(http.StatusFound, s.Config.ForwardedPrefix)
	}
}

func (s *Server) handleGetDeletionPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		short := c.Param("short")
		c.HTML(http.StatusOK, "delete.gohtml", short)
	}
}

func (s *Server) handlePostDeletionPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		short := c.Param("short")
		if err := s.Shortener.DeleteShortURL(short); err != nil {
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
		user := session.Get("user_id")
		if user == nil {
			c.Redirect(http.StatusFound, loginPage)
			return
		}
		c.Next()
	}
}
