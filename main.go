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

	pool := openPGConnectionPool(err, conf)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	accounts := gin.Accounts{conf.BasicAuthUser: conf.BasicAuthSecret}

	redirectsRepository := repository.NewInMemoryRedirectRepository()
	usersRepository := repository.NewDBUsersRepository(pool)
	urlShortener := shortener.NewUrlShortenerService(hasher.UrlHasher, redirectsRepository, usersRepository)

	r := gin.Default()
	r = handler.SetupRoutes(r, conf.ForwardedPrefix, urlShortener, accounts)

	if err := r.Run(); err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}

func openPGConnectionPool(err error, config config.Configuration) *pgxpool.Pool {
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
