package main

import (
	stdlog "log"
	"net/http"

	"goworld-skeleton/internal/config"
	"goworld-skeleton/internal/dao"
	logger "goworld-skeleton/internal/log"
	"goworld-skeleton/internal/modules/account"
	"goworld-skeleton/internal/modules/bag"
	"goworld-skeleton/internal/modules/chat"
	"goworld-skeleton/internal/modules/item"
	"goworld-skeleton/internal/modules/mail"
	"goworld-skeleton/internal/modules/match"
	"goworld-skeleton/internal/modules/notice"
	"goworld-skeleton/internal/modules/player"
	"goworld-skeleton/internal/modules/room"
	"goworld-skeleton/internal/modules/shop"
	"goworld-skeleton/internal/redis"
	"goworld-skeleton/internal/server"
)

func main() {
	cfg := config.Default()
	log := logger.New(cfg.Environment)
	store := dao.NewDataStore()
	cache := redis.NewCache()

	services := server.Services{
		Account: account.NewService(store, cache, log),
		Player:  player.NewService(store, log),
		Bag:     bag.NewService(store, log),
		Item:    item.NewService(store),
		Shop:    shop.NewService(store, log),
		Mail:    mail.NewService(store, log),
		Notice:  notice.NewService(store),
		Chat:    chat.NewService(store, log),
		Room:    room.NewService(store, log),
		Match:   match.NewService(store, log),
	}

	handler := server.NewRouter(services)
	log.Printf("GoWorld skeleton server listening on %s", cfg.HTTPPort)
	if err := http.ListenAndServe(cfg.HTTPPort, handler); err != nil {
		stdlog.Fatalf("failed to start server: %v", err)
	}
}
