package internal

import "github.com/gin-gonic/gin"

type RedirectCreationRequest struct {
	Url string `json:"url"`
}

func CreateHealthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Status(200)
		return
	}
}

func CreateGetHandler(shortener UrlShortener) gin.HandlerFunc {
	return func(c *gin.Context) {
		short := c.Param("short")

		expand, err := shortener.Expand(short)
		if err != nil {
			c.JSON(404, gin.H{"error": "Redirect not found"})
		}

		c.Redirect(302, expand)
	}
}

func CreatePostHandler(shortener UrlShortener) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request RedirectCreationRequest
		err := c.BindJSON(&request)
		if err != nil {
			return
		}

		if request.Url == "" {
			c.JSON(400, gin.H{"error": "Url is required"})
			return
		}

		short, err := shortener.Shorten(request.Url)
		if err != nil {
			c.JSON(500, gin.H{"error": "Internal server error"})
			return
		}

		c.JSON(201, gin.H{"short": short})
	}
}

func CreateDeleteHandler(shortener UrlShortener) gin.HandlerFunc {
	return func(c *gin.Context) {
		short := c.Param("short")

		err := shortener.Delete(short)
		if err != nil {
			c.JSON(404, gin.H{"error": "Redirect not found"})
		}

		c.JSON(200, gin.H{"success": "Redirect deleted"})
	}
}
