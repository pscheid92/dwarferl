package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
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

	pool := openPGConnectionPool(err, config)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	accounts := gin.Accounts{config.BasicAuthUser: config.BasicAuthSecret}

	redirectsRepository := repository.NewInMemoryRedirectRepository()
	usersRepository := repository.NewDBUsersRepository(pool)
	urlShortener := shortener.NewUrlShortenerService(hasher.UrlHasher, redirectsRepository, usersRepository)

	r := gin.Default()
	r = handler.SetupRoutes(r, config.ForwardedPrefix, urlShortener, accounts)

	if err := r.Run(); err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}

func openPGConnectionPool(err error, config internal.Configuration) *pgxpool.Pool {
	c, err := pgxpool.ParseConfig("")
	if err != nil {
		log.Fatal(err)
	}

	c.ConnConfig.Host = config.Database.Host
	c.ConnConfig.Port = config.Database.Port
	c.ConnConfig.Database = config.Database.Name
	c.ConnConfig.User = config.Database.User
	c.ConnConfig.Password = config.Database.Password

	pool, err := pgxpool.ConnectConfig(context.Background(), c)
	return pool
}
