package bag

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"goworld-skeleton/internal/dao"
)

// Service exposes bag operations.
type Service struct {
	store  *dao.DataStore
	logger *log.Logger
}

// NewService constructs a bag service.
func NewService(store *dao.DataStore, logger *log.Logger) Service {
	return Service{store: store, logger: logger}
}

// Register binds HTTP endpoints.
func (s Service) Register(mux *http.ServeMux) {
	mux.HandleFunc("/api/bag/", s.getBag)
}

func (s Service) getBag(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	playerID := strings.TrimPrefix(r.URL.Path, "/api/bag/")
	var bag []dao.BagEntry
	s.store.WithRead(func(store *dao.DataStore) {
		bag = store.Bags[playerID]
	})

	if bag == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "bag not found"})
		return
	}

	s.logger.Printf("bag fetched for %s", playerID)
	writeJSON(w, http.StatusOK, map[string]interface{}{"items": bag})
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
