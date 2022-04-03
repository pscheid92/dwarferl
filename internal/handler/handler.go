package handler

import (
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
		authorized.GET("/", createListHandler(shortener))
		authorized.POST("/", createPostHandler(shortener))
		authorized.DELETE("/:short", createDeleteHandler(shortener))
	}

	return router
}

type RedirectCreationRequest struct {
	Url string `json:"url"`
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
		c.Redirect(http.StatusMovedPermanently, redirect.URL)
	}
}

func createListHandler(shortener internal.UrlShortenerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		list, err := shortener.List("00000000-0000-0000-0000-000000000000")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, list)
	}
}

func createPostHandler(shortener internal.UrlShortenerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request RedirectCreationRequest
		err := c.BindJSON(&request)
		if err != nil {
			return
		}

		if request.Url == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Url is required"})
			return
		}

		short, err := shortener.ShortenURL(request.Url)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error", "details": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"short": short})
	}
}

func createDeleteHandler(shortener internal.UrlShortenerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		short := c.Param("short")

		err := shortener.DeleteShortURL(short)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Redirect not found", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": "Redirect deleted"})
	}
}
