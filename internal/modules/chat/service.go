package chat

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"goworld-skeleton/internal/dao"
)

// Service exposes chat operations.
type Service struct {
	store  *dao.DataStore
	logger *log.Logger
}

// NewService constructs a chat service.
func NewService(store *dao.DataStore, logger *log.Logger) Service {
	return Service{store: store, logger: logger}
}

// Register binds HTTP endpoints.
func (s Service) Register(mux *http.ServeMux) {
	mux.HandleFunc("/api/chat/", s.handler)
	mux.HandleFunc("/api/chat", s.handler)
}

type messageInput struct {
	From    string `json:"from"`
	To      string `json:"to"`
	RoomID  string `json:"room_id"`
	Body    string `json:"body"`
	Channel string `json:"channel"`
}

func (s Service) handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.send(w, r)
	case http.MethodGet:
		s.history(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (s Service) send(w http.ResponseWriter, r *http.Request) {
	var input messageInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	msg := dao.ChatMessage{
		From:    input.From,
		To:      input.To,
		RoomID:  input.RoomID,
		Body:    input.Body,
		Channel: input.Channel,
		SentAt:  time.Now(),
	}

	s.store.WithLock(func(store *dao.DataStore) {
		store.Chats = append(store.Chats, msg)
	})

	s.logger.Printf("chat message recorded from %s", msg.From)
	writeJSON(w, http.StatusCreated, msg)
}

func (s Service) history(w http.ResponseWriter, _ *http.Request) {
	var history []dao.ChatMessage
	s.store.WithRead(func(store *dao.DataStore) { history = append(history, store.Chats...) })
	writeJSON(w, http.StatusOK, map[string]interface{}{"messages": history})
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
