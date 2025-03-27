package main

import (
	"flag"
	"log"

	"github.com/escalopa/peer-cast/internal/app"
	"github.com/escalopa/peer-cast/internal/auth"
	"github.com/escalopa/peer-cast/internal/config"
	"github.com/escalopa/peer-cast/internal/db"
	"github.com/escalopa/peer-cast/internal/room"
	"github.com/escalopa/peer-cast/internal/service"
)

var configPath = flag.String("config", "config.yml", "path to config file")

func main() {
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	database, err := db.New(cfg.DB.File)
	if err != nil {
		log.Fatalf("init database: %v", err)
	}
	defer func() { _ = database.Close() }()

	hubInstance := room.NewHub()
	userTokenProvider := auth.NewUserProvider(cfg.JWT.User)
	roomTokenProvider := auth.NewRoomProvider(cfg.JWT.Room)
	oauthProvider := auth.NewOAuthProvider(cfg.OAuth)

	s := app.New(
		app.Config{
			Domain:          cfg.App.Domain,
			AllowOrigins:    cfg.App.AllowOrigins,
			AccessTokenTTL:  cfg.JWT.User.AccessTokenTTL,
			RefreshTokenTTL: cfg.JWT.User.RefreshTokenTTL,
		}, service.New(
			database, hubInstance, oauthProvider, userTokenProvider, roomTokenProvider,
		),
	)

	if err := s.Run(cfg.App.Addr); err != nil {
		log.Fatalf("server start: %v", err)
	}
}
