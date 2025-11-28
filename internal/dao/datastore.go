package dao

import (
	"sync"
	"time"
)

// DataStore is an in-memory stand-in for a relational database.
type DataStore struct {
	mu sync.RWMutex

	Accounts map[string]Account
	Players  map[string]Player
	Items    []Item
	Notices  []Notice
	Mails    map[string][]Mail
	Bags     map[string][]BagEntry
	Chats    []ChatMessage
	Rooms    map[string]Room
}

// NewDataStore seeds a datastore with demo data.
func NewDataStore() *DataStore {
	items := []Item{
		{ID: "potion", Name: "Small Potion", Rarity: "common", Price: 25},
		{ID: "sword", Name: "Bronze Sword", Rarity: "uncommon", Price: 120},
	}

	notices := []Notice{{ID: "welcome", Title: "Welcome", Body: "服务器已启动，祝你游戏愉快！", Severity: "info", CreatedAt: time.Now()}}

	players := map[string]Player{
		"demo": {ID: "demo", Name: "DemoPlayer", Level: 10, Experience: 2200, LastLogin: time.Now()},
	}

	return &DataStore{
		Accounts: map[string]Account{"demo": {ID: "demo", Username: "demo", Password: "password", Token: "demo-token"}},
		Players:  players,
		Items:    items,
		Notices:  notices,
		Mails:    map[string][]Mail{"demo": {{ID: "m1", Subject: "欢迎礼包", Body: "感谢试玩", Attachments: []MailAttachment{{ItemID: "potion", Quantity: 3}}}}},
		Bags:     map[string][]BagEntry{"demo": {{ItemID: "potion", Quantity: 2}}},
		Chats:    []ChatMessage{},
		Rooms:    map[string]Room{},
	}
}

func (d *DataStore) WithLock(fn func(store *DataStore)) {
	d.mu.Lock()
	defer d.mu.Unlock()
	fn(d)
}

func (d *DataStore) WithRead(fn func(store *DataStore)) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	fn(d)
}

// Domain models to be shared across modules.
type Account struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
	Token    string `json:"token"`
}

type Player struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Level      int       `json:"level"`
	Experience int       `json:"experience"`
	LastLogin  time.Time `json:"last_login"`
}

type BagEntry struct {
	ItemID   string `json:"item_id"`
	Quantity int    `json:"quantity"`
}

type Item struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Rarity string `json:"rarity"`
	Price  int    `json:"price"`
}

type ShopListing struct {
	ItemID string `json:"item_id"`
	Price  int    `json:"price"`
	Stock  int    `json:"stock"`
}

type MailAttachment struct {
	ItemID   string `json:"item_id"`
	Quantity int    `json:"quantity"`
}

type Mail struct {
	ID          string           `json:"id"`
	Subject     string           `json:"subject"`
	Body        string           `json:"body"`
	Attachments []MailAttachment `json:"attachments"`
}

type Notice struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Severity  string    `json:"severity"`
	CreatedAt time.Time `json:"created_at"`
}

type ChatMessage struct {
	From    string    `json:"from"`
	To      string    `json:"to,omitempty"`
	RoomID  string    `json:"room_id,omitempty"`
	Body    string    `json:"body"`
	SentAt  time.Time `json:"sent_at"`
	Channel string    `json:"channel"`
}

type Room struct {
	ID         string   `json:"id"`
	Game       string   `json:"game"`
	Players    []string `json:"players"`
	MaxPlayers int      `json:"max_players"`
	Status     string   `json:"status"`
}
