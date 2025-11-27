package shop

import (
	"encoding/json"
	"log"
	"net/http"

	"goworld-skeleton/internal/dao"
)

// Service exposes simple shop operations.
type Service struct {
	store  *dao.DataStore
	logger *log.Logger
}

// NewService constructs a shop service.
func NewService(store *dao.DataStore, logger *log.Logger) Service {
	return Service{store: store, logger: logger}
}

// Register binds HTTP endpoints.
func (s Service) Register(mux *http.ServeMux) {
	mux.HandleFunc("/api/shop/items", s.list)
}

func (s Service) list(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	listings := make([]dao.ShopListing, 0)
	s.store.WithRead(func(store *dao.DataStore) {
		for _, item := range store.Items {
			listings = append(listings, dao.ShopListing{ItemID: item.ID, Price: item.Price, Stock: 99})
		}
	})

	s.logger.Printf("shop listing served")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"items": listings})
}
