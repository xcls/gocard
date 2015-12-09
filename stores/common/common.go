package common

import (
	"time"
)

type Store struct {
	Cards CardStore
	Decks DeckStore
	Users UserStore
}

type Card struct {
	ID        int64
	Context   string
	Front     string
	Back      string
	DeckID    int64
	CreatedAt time.Time
}

type Deck struct {
	ID        int64
	Name      string
	CreatedAt time.Time
	// UpdatedAt time.Time
}

type User struct {
	ID                int64
	Email             string
	EncryptedPassword string `json:"-"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type CardStore interface {
	Insert(model *Card) error
	Update(model *Card) error
	All() ([]*Card, error)
	AllByDeckID(id int) ([]*Card, error)
	Find(id int64) (*Card, error)
}

type DeckStore interface {
	Insert(deck *Deck) error
	Find(id int64) (*Deck, error)
	All() ([]*Deck, error)
}

type UserStore interface {
	Insert(model *User) error
	Find(id int64) (*User, error)
	Authenticate(email, password string) (*User, error)
}
