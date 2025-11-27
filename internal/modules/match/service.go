package match

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"goworld-skeleton/internal/dao"
)

// Service exposes matchmaking endpoints.
type Service struct {
	store  *dao.DataStore
	logger *log.Logger
}

// NewService constructs a match service.
func NewService(store *dao.DataStore, logger *log.Logger) Service {
	return Service{store: store, logger: logger}
}

// Register binds HTTP endpoints.
func (s Service) Register(mux *http.ServeMux) {
	mux.HandleFunc("/api/match/enqueue", s.enqueue)
}

type enqueueInput struct {
	PlayerID string `json:"player_id"`
	Mode     string `json:"mode"`
}

func (s Service) enqueue(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}

	var input enqueueInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	matchID := generateMatchID()
	room := dao.Room{ID: matchID, Game: input.Mode, Players: []string{input.PlayerID}, MaxPlayers: 4, Status: "matching"}
	s.store.WithLock(func(store *dao.DataStore) { store.Rooms[room.ID] = room })

	s.logger.Printf("player %s enqueued for %s", input.PlayerID, input.Mode)
	writeJSON(w, http.StatusAccepted, map[string]interface{}{"match_id": matchID, "room": room})
}

func generateMatchID() string {
	return "match-" + time.Now().Format("150405.000")
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
