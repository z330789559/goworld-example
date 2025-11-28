package mail

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"goworld-skeleton/internal/dao"
)

// Service exposes mail endpoints.
type Service struct {
	store  *dao.DataStore
	logger *log.Logger
}

// NewService constructs a mail service.
func NewService(store *dao.DataStore, logger *log.Logger) Service {
	return Service{store: store, logger: logger}
}

// Register binds HTTP endpoints.
func (s Service) Register(mux *http.ServeMux) {
	mux.HandleFunc("/api/mail/", s.list)
}

func (s Service) list(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	playerID := strings.TrimPrefix(r.URL.Path, "/api/mail/")
	var mails []dao.Mail
	s.store.WithRead(func(store *dao.DataStore) {
		mails = store.Mails[playerID]
	})

	if mails == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "no mail for player"})
		return
	}

	s.logger.Printf("mail fetched for %s", playerID)
	writeJSON(w, http.StatusOK, map[string]interface{}{"messages": mails})
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
