package item

import (
	"encoding/json"
	"net/http"

	"goworld-skeleton/internal/dao"
)

// Service exposes item catalog read-only operations.
type Service struct {
	store *dao.DataStore
}

// NewService constructs an item service.
func NewService(store *dao.DataStore) Service {
	return Service{store: store}
}

// Register binds HTTP endpoints.
func (s Service) Register(mux *http.ServeMux) {
	mux.HandleFunc("/api/items/", s.list)
	mux.HandleFunc("/api/items", s.list)
}

func (s Service) list(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}
	var items []dao.Item
	s.store.WithRead(func(store *dao.DataStore) { items = store.Items })

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"items": items})
}
