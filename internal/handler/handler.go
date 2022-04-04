package handler

import (
	"errors"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/pscheid92/dwarferl/internal"
	"net/http"
)

func SetupRoutes(router *gin.Engine, forwardedPrefix string, shortener internal.UrlShortenerService, cookies sessions.CookieStore) *gin.Engine {
	router.RedirectTrailingSlash = false

	router.Use(sessions.Sessions("dwarferl_session", cookies))

	g := router.Group(forwardedPrefix)
	g.GET("/health", createHealthHandler())
	g.GET("/:short", createGetHandler(shortener))

	g.GET("/login", createLoginPage(forwardedPrefix))
	g.POST("/login", handleLogin(forwardedPrefix))

	g.GET("/logout", handleLogout(forwardedPrefix))

	authorized := g.Group("")
	authorized.Use(authRequired(forwardedPrefix + "login"))
	{
		authorized.GET("/", indexPage(shortener, forwardedPrefix))
		authorized.GET("/create", serveCreationPage())
		authorized.POST("/create", handleCreationPage(shortener, forwardedPrefix))
		authorized.GET("/delete/:short", serverDeletionPage())
		authorized.POST("/delete/:short", handleDeletionPage(shortener, forwardedPrefix))
	}

	return router
}

type RedirectCreationRequest struct {
	Url string `form:"url"`
}

func createLoginPage(linkPrefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")

		data := gin.H{
			"user_id":     userID,
			"link_prefix": linkPrefix,
		}

		c.HTML(http.StatusOK, "login.gohtml", data)
	}
}

func handleLogin(linkPrefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		username := c.PostForm("username")
		password := c.PostForm("password")

		if username != "admin" || password != "admin" {
			c.Redirect(http.StatusFound, linkPrefix)
			return
		}

		session.Set("user_id", "00000000-0000-0000-0000-000000000000")
		if err := session.Save(); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.Redirect(http.StatusFound, linkPrefix)
	}
}

func handleLogout(linkPrefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user_id")
		if user == nil {
			c.Redirect(http.StatusBadRequest, linkPrefix+"login")
			return
		}

		session.Delete("user_id")
		if err := session.Save(); err != nil {
			c.Redirect(http.StatusInternalServerError, linkPrefix+"login")
			return
		}

		c.Redirect(http.StatusFound, linkPrefix+"login")
	}
}

func authRequired(redirectPage string) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user_id")
		if user == nil {
			c.Redirect(http.StatusTemporaryRedirect, redirectPage)
			return
		}
		c.Next()
	}
}

func createHealthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Status(http.StatusOK)
		return
	}
}

func createGetHandler(shortener internal.UrlShortenerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		short := c.Param("short")

		redirect, err := shortener.ExpandShortURL(short)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Redirect not found", "details": err.Error()})
			return
		}

		c.Header("Cache-Control", "private, max-age=90")
		c.Header("Referrer-Policy", "unsafe-url")
		c.Redirect(http.StatusMovedPermanently, redirect)
	}
}

func indexPage(shortener internal.UrlShortenerService, linkPrefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		list, err := shortener.List("00000000-0000-0000-0000-000000000000")
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		data := gin.H{
			"redirects":  list,
			"linkPrefix": linkPrefix,
		}

		c.HTML(http.StatusOK, "index.gohtml", data)
	}
}

func serveCreationPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "create.gohtml", nil)
	}
}

func handleCreationPage(shortener internal.UrlShortenerService, linkPrefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request RedirectCreationRequest

		err := c.Bind(&request)
		if err != nil {
			return
		}

		if request.Url == "" {
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("Url is required"))
			return
		}

		_, err = shortener.ShortenURL(request.Url)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.Redirect(http.StatusFound, linkPrefix)
	}
}

func serverDeletionPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		short := c.Param("short")
		c.HTML(http.StatusOK, "delete.gohtml", short)
	}
}

func handleDeletionPage(shortener internal.UrlShortenerService, linkPrefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		short := c.Param("short")
		if err := shortener.DeleteShortURL(short); err != nil {
			_ = c.AbortWithError(http.StatusNotFound, err)
			return
		}
		c.Redirect(http.StatusFound, linkPrefix)
	}
}
