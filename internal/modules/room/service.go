package room

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"goworld-skeleton/internal/dao"
)

// Service exposes lightweight room orchestration.
type Service struct {
	store  *dao.DataStore
	logger *log.Logger
}

// NewService constructs a room service.
func NewService(store *dao.DataStore, logger *log.Logger) Service {
	return Service{store: store, logger: logger}
}

// Register binds HTTP endpoints.
func (s Service) Register(mux *http.ServeMux) {
	mux.HandleFunc("/api/room/create", s.create)
	mux.HandleFunc("/api/room/", s.list)
	mux.HandleFunc("/api/room", s.list)
}

type createInput struct {
	Game       string   `json:"game"`
	Players    []string `json:"players"`
	MaxPlayers int      `json:"max_players"`
}

func (s Service) create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}

	var input createInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	room := dao.Room{ID: generateRoomID(), Game: input.Game, Players: input.Players, MaxPlayers: input.MaxPlayers, Status: "waiting"}
	s.store.WithLock(func(store *dao.DataStore) { store.Rooms[room.ID] = room })

	s.logger.Printf("room %s created for %s", room.ID, room.Game)
	writeJSON(w, http.StatusCreated, room)
}

func (s Service) list(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}
	rooms := make([]dao.Room, 0)
	s.store.WithRead(func(store *dao.DataStore) {
		for _, r := range store.Rooms {
			rooms = append(rooms, r)
		}
	})
	writeJSON(w, http.StatusOK, map[string]interface{}{"rooms": rooms})
}

func generateRoomID() string {
	return "room-" + time.Now().Format("150405.000")
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
