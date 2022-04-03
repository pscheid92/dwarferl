package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pscheid92/dwarferl/internal/config"
	"github.com/pscheid92/dwarferl/internal/handler"
	"github.com/pscheid92/dwarferl/internal/hasher"
	"github.com/pscheid92/dwarferl/internal/repository"
	"github.com/pscheid92/dwarferl/internal/shortener"
	"log"
)

func main() {
	conf, err := config.GatherConfig()
	if err != nil {
		log.Fatal(err)
	}

	pool, err := openPGConnectionPool(conf)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	accounts := gin.Accounts{conf.BasicAuthUser: conf.BasicAuthSecret}

	redirectsRepository := repository.NewDBRedirectsRepository(pool)
	usersRepository := repository.NewDBUsersRepository(pool)
	urlShortener := shortener.NewUrlShortenerService(hasher.UrlHasher, redirectsRepository, usersRepository)

	r := gin.Default()
	r = handler.SetupRoutes(r, conf.ForwardedPrefix, urlShortener, accounts)

	if err := r.Run(); err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}

func openPGConnectionPool(config config.Configuration) (*pgxpool.Pool, error) {
	c, err := pgxpool.ParseConfig(config.DatabaseURL)
	if err != nil {
		return nil, err
	}
	c.ConnConfig.Password = config.DatabasePassword
	return pgxpool.ConnectConfig(context.Background(), c)
}
