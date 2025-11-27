package notice

import (
	"encoding/json"
	"net/http"

	"goworld-skeleton/internal/dao"
)

// Service exposes notice board endpoints.
type Service struct {
	store *dao.DataStore
}

// NewService constructs a notice service.
func NewService(store *dao.DataStore) Service {
	return Service{store: store}
}

// Register binds HTTP endpoints.
func (s Service) Register(mux *http.ServeMux) {
	mux.HandleFunc("/api/notice/", s.list)
	mux.HandleFunc("/api/notice", s.list)
}

func (s Service) list(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}
	var notices []dao.Notice
	s.store.WithRead(func(store *dao.DataStore) { notices = store.Notices })

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"notices": notices})
}
