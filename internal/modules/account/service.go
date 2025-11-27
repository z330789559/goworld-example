package account

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"goworld-skeleton/internal/dao"
	cache "goworld-skeleton/internal/redis"
)

// Service exposes account use cases.
type Service struct {
	store  *dao.DataStore
	cache  *cache.Cache
	logger *log.Logger
}

// NewService constructs an account service.
func NewService(store *dao.DataStore, cache *cache.Cache, logger *log.Logger) Service {
	return Service{store: store, cache: cache, logger: logger}
}

// Register registers HTTP handlers on the provided mux.
func (s Service) Register(mux *http.ServeMux) {
	mux.HandleFunc("/api/account/register", s.register)
	mux.HandleFunc("/api/account/login", s.login)
}

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s Service) register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}

	var input credentials
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if input.Username == "" || input.Password == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "username and password required"})
		return
	}

	var created dao.Account
	s.store.WithLock(func(store *dao.DataStore) {
		if _, exists := store.Accounts[input.Username]; exists {
			return
		}
		token := generateToken()
		created = dao.Account{ID: input.Username, Username: input.Username, Password: input.Password, Token: token}
		store.Accounts[input.Username] = created
		store.Players[input.Username] = dao.Player{ID: input.Username, Name: input.Username, Level: 1, Experience: 0, LastLogin: time.Now()}
	})

	if created.ID == "" {
		writeJSON(w, http.StatusConflict, map[string]string{"error": "username already exists"})
		return
	}

	s.logger.Printf("registered user %s", created.Username)
	writeJSON(w, http.StatusCreated, map[string]string{"token": created.Token})
}

func (s Service) login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}

	var input credentials
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	account, err := s.authenticate(input.Username, input.Password)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	s.cache.Set("session:"+account.ID, account.Token, 30*time.Minute)
	s.logger.Printf("user %s logged in", account.Username)
	writeJSON(w, http.StatusOK, map[string]string{"token": account.Token})
}

func (s Service) authenticate(username, password string) (dao.Account, error) {
	var account dao.Account
	s.store.WithRead(func(store *dao.DataStore) {
		if acc, ok := store.Accounts[username]; ok && acc.Password == password {
			account = acc
		}
	})

	if account.ID == "" {
		return dao.Account{}, errors.New("invalid credentials")
	}

	return account, nil
}

func generateToken() string {
	return time.Now().Format("20060102150405.000000000")
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
