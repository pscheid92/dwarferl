package main

import (
	"github.com/gin-gonic/gin"
	"github.com/pscheid92/dwarferl/internal/handler"
	"github.com/pscheid92/dwarferl/internal/hasher"
	"github.com/pscheid92/dwarferl/internal/repository"
	"github.com/pscheid92/dwarferl/internal/shortener"
	"log"
	"os"
	"strings"
)

func main() {
	accounts := getBasicAuthAccounts()
	forwardedPrefix := prepareForwardedPrefix()

	redirectsRepository := repository.NewInMemoryRedirectRepository()
	usersRepository := repository.StaticUsersRepository{}
	urlShortener := shortener.NewUrlShortenerService(hasher.UrlHasher, redirectsRepository, usersRepository)

	r := gin.Default()
	r = handler.SetupRoutes(gin.Default(), forwardedPrefix, urlShortener, accounts)

	if err := r.Run(); err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}

func prepareForwardedPrefix() string {
	forwardedPrefix := os.Getenv("FORWARDED_PREFIX")
	if !strings.HasSuffix(forwardedPrefix, "/") {
		forwardedPrefix += "/"
	}
	return forwardedPrefix
}

func getBasicAuthAccounts() gin.Accounts {
	user := os.Getenv("DWARFERL_USER")
	secret := os.Getenv("DWARFERL_SECRET")
	if user == "" || secret == "" {
		log.Fatal("DWARFERL_USER and DWARFERL_SECRET must be set")
	}
	return gin.Accounts{user: secret}
}
