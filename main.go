package main

import (
	"github.com/gin-gonic/gin"
	"github.com/pscheid92/dwarferl/internal"
	"github.com/pscheid92/dwarferl/internal/handler"
	"github.com/pscheid92/dwarferl/internal/hasher"
	"github.com/pscheid92/dwarferl/internal/repository"
	"github.com/pscheid92/dwarferl/internal/shortener"
	"log"
)

func main() {
	config, err := internal.GatherConfig()
	if err != nil {
		log.Fatal(err)
	}

	accounts := gin.Accounts{config.BasicAuthUser: config.BasicAuthSecret}

	redirectsRepository := repository.NewInMemoryRedirectRepository()
	usersRepository := repository.StaticUsersRepository{}
	urlShortener := shortener.NewUrlShortenerService(hasher.UrlHasher, redirectsRepository, usersRepository)

	r := gin.Default()
	r = handler.SetupRoutes(r, config.ForwardedPrefix, urlShortener, accounts)

	if err := r.Run(); err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}
