package main

import (
	"context"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/markbates/goth/gothic"
	"github.com/pscheid92/dwarferl/internal/config"
	"github.com/pscheid92/dwarferl/internal/hasher"
	"github.com/pscheid92/dwarferl/internal/repository"
	"github.com/pscheid92/dwarferl/internal/server"
	"github.com/pscheid92/dwarferl/internal/shortener"
	"log"
)

func main() {
	conf, err := config.GatherConfig()
	if err != nil {
		log.Fatal(err)
	}

	pool, err := openPGConnectionPool()
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	sessionStore := cookie.NewStore([]byte(conf.SessionSecret))
	gothic.Store = cookie.NewStore([]byte(conf.SessionSecret))

	redirectsRepository := repository.NewDBRedirectsRepository(pool)
	usersRepository := repository.NewDBUsersRepository(pool)
	urlShortener := shortener.NewUrlShortenerService(hasher.UrlHasher, redirectsRepository, usersRepository)

	svr := server.New(conf, sessionStore, urlShortener)
	svr.Use(gin.Logger(), gin.Recovery())
	svr.InitRoutes()

	if err := svr.Run(); err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}

func openPGConnectionPool() (*pgxpool.Pool, error) {
	c, err := pgxpool.ParseConfig("")
	if err != nil {
		return nil, err
	}
	return pgxpool.ConnectConfig(context.Background(), c)
}
