package common

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var ErrUserAuthFailed = errors.New("Email and password don't match or user doesn't exist.")

const (
	UserPasswordCost = 11
)

type Store struct {
	Cards        CardStore
	Decks        DeckStore
	Users        UserStore
	UserSessions UserSessionStore
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

type UserSession struct {
	ID        int64
	UID       string `json:"-"`
	UserID    int64
	CreatedAt time.Time
}

func NewUserSession(userID int64) *UserSession {
	return &UserSession{
		UID:    generateSessionID(),
		UserID: userID,
	}
}

// generateSessionID generates a random session ID
func generateSessionID() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
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

type UserSessionStore interface {
	Insert(model *UserSession) error
	Find(uid string) (*UserSession, error)
}
