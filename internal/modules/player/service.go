package player

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"goworld-skeleton/internal/dao"
)

// Service exposes player endpoints.
type Service struct {
	store  *dao.DataStore
	logger *log.Logger
}

// NewService constructs a player service.
func NewService(store *dao.DataStore, logger *log.Logger) Service {
	return Service{store: store, logger: logger}
}

// Register binds HTTP endpoints.
func (s Service) Register(mux *http.ServeMux) {
	mux.HandleFunc("/api/player/", s.getProfile)
}

func (s Service) getProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	playerID := strings.TrimPrefix(r.URL.Path, "/api/player/")
	var player dao.Player
	s.store.WithRead(func(store *dao.DataStore) {
		player = store.Players[playerID]
	})

	if player.ID == "" {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "player not found"})
		return
	}

	s.logger.Printf("loaded player %s", player.ID)
	writeJSON(w, http.StatusOK, player)
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
