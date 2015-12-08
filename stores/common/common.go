package common

import (
	"time"
)

type Store struct {
	Cards CardStore
	Decks DeckStore
}

type Card struct {
	ID        int64
	Context   string
	Front     string
	Back      string
	DeckID    int64
	CreatedAt time.Time
}

type CardStore interface {
	Insert(model *Card) error
	Update(model *Card) error
	All() ([]*Card, error)
	AllByDeckID(id int) ([]*Card, error)
	Find(id int64) (*Card, error)
}

type Deck struct {
	ID        int64
	Name      string
	CreatedAt time.Time
	// UpdatedAt time.Time
}

type DeckStore interface {
	Insert(deck *Deck) error
	Find(id int64) (*Deck, error)
	All() ([]*Deck, error)
}
