package common

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	UserPasswordCost = 11
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

func (m *User) SetPassword(pass string) error {
	encrypted, err := bcrypt.GenerateFromPassword([]byte(pass), UserPasswordCost)
	if err != nil {
		return err
	}
	m.EncryptedPassword = string(encrypted)
	return nil
}

// ComparePassword returns nil if equal, error if not
func (m *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(m.EncryptedPassword),
		[]byte(password),
	)
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
	FindByEmail(email string) (*User, error)
}
