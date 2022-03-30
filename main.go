package main

import (
	"github.com/gin-gonic/gin"
	"github.com/pscheid92/dwarferl/internal"
	"log"
	"os"
	"strings"
)

func main() {
	accounts := getBasicAuthAccounts()
	forwardedPrefix := prepareForwardedPrefix()

	redirectsRepository := internal.NewInMemoryRedirectRepository()
	usersRepository := internal.StaticUsersRepository{}
	shortener := internal.NewUrlShortenerService(internal.UrlHasher, redirectsRepository, usersRepository)

	r := gin.Default()
	r = internal.SetupRoutes(gin.Default(), forwardedPrefix, shortener, accounts)

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
