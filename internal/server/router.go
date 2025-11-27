package server

import (
	"encoding/json"
	"net/http"
)

// Services aggregates domain services for dependency injection.
type Services struct {
	Account AccountRoutes
	Player  PlayerRoutes
	Bag     BagRoutes
	Item    ItemRoutes
	Shop    ShopRoutes
	Mail    MailRoutes
	Notice  NoticeRoutes
	Chat    ChatRoutes
	Room    RoomRoutes
	Match   MatchRoutes
}

// NewRouter wires HTTP handlers for all modules.
func NewRouter(services Services) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	services.Account.Register(mux)
	services.Player.Register(mux)
	services.Bag.Register(mux)
	services.Item.Register(mux)
	services.Shop.Register(mux)
	services.Mail.Register(mux)
	services.Notice.Register(mux)
	services.Chat.Register(mux)
	services.Room.Register(mux)
	services.Match.Register(mux)

	return mux
}

// Shared helpers and interfaces for handlers.
type AccountRoutes interface{ Register(*http.ServeMux) }
type PlayerRoutes interface{ Register(*http.ServeMux) }
type BagRoutes interface{ Register(*http.ServeMux) }
type ItemRoutes interface{ Register(*http.ServeMux) }
type ShopRoutes interface{ Register(*http.ServeMux) }
type MailRoutes interface{ Register(*http.ServeMux) }
type NoticeRoutes interface{ Register(*http.ServeMux) }
type ChatRoutes interface{ Register(*http.ServeMux) }
type RoomRoutes interface{ Register(*http.ServeMux) }
type MatchRoutes interface{ Register(*http.ServeMux) }

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
