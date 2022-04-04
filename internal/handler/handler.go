package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/pscheid92/dwarferl/internal"
	"net/http"
)

func SetupRoutes(router *gin.Engine, forwardedPrefix string, shortener internal.UrlShortenerService, accounts gin.Accounts) *gin.Engine {
	router.RedirectTrailingSlash = false

	g := router.Group(forwardedPrefix)
	g.GET("/health", createHealthHandler())
	g.GET("/:short", createGetHandler(shortener))

	authorized := g.Group("", gin.BasicAuth(accounts))
	{
		authorized.GET("/", indexPage(shortener, forwardedPrefix))
		authorized.GET("/create", serveCreationPage())
		authorized.POST("/create", handleCreationPage(shortener))
		authorized.GET("/delete/:short", serverDeletionPage())
		authorized.POST("/delete/:short", handleDeletionPage(shortener))
	}

	return router
}

type RedirectCreationRequest struct {
	Url string `form:"url"`
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

func handleCreationPage(shortener internal.UrlShortenerService) gin.HandlerFunc {
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

		c.Redirect(http.StatusFound, "/")
	}
}

func serverDeletionPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		short := c.Param("short")
		c.HTML(http.StatusOK, "delete.gohtml", short)
	}
}

func handleDeletionPage(shortener internal.UrlShortenerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		short := c.Param("short")
		if err := shortener.DeleteShortURL(short); err != nil {
			_ = c.AbortWithError(http.StatusNotFound, err)
			return
		}
		c.Redirect(http.StatusFound, "/")
	}
}
