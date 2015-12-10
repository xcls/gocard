package common

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"math"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var ErrUserAuthFailed = errors.New("Email and password don't match or user doesn't exist.")

const (
	UserPasswordCost = 11
	MinEF            = 1.3
	MaxEF            = 2.5
)

type Store struct {
	CardReviews  CardReviewStore
	Cards        CardStore
	Decks        DeckStore
	Reviews      ReviewStore
	UserSessions UserSessionStore
	Users        UserStore
}

type Card struct {
	ID        int64
	Context   string
	Front     string
	Back      string
	DeckID    int64
	CreatedAt time.Time
}

// CardReview is a combination of Review, Card and Deck so it can be
// conveniently retrieved from stores
type CardReview struct {
	ID         int64
	Enabled    bool
	EaseFactor float64
	Interval   int64
	DueOn      time.Time
	UserID     int64

	CardID      int64
	CardContext string
	CardFront   string
	CardBack    string

	DeckID   int64
	DeckName string
}

type Deck struct {
	ID        int64
	Name      string
	CreatedAt time.Time
	// UpdatedAt time.Time
}

type Review struct {
	ID         int64
	Enabled    bool      // Enabled for review or not?
	EaseFactor float64   // The E-Factor is used to determine next interval length
	Interval   int64     // Indicates the last interval in days
	DueOn      time.Time // When the next review is due
	CardID     int64
	UserID     int64
	CreatedAt  time.Time
}

func (r *Review) AddRating(rating int64) error {
	if rating < 0 || rating > 5 {
		return errors.New("rating must be between 0 and 5")
	}
	r.EaseFactor = UpdateEF(r.EaseFactor, rating)
	r.Interval = UpdateInterval(r.Interval, r.EaseFactor)
	return nil
}

// UpdateEF calculates the new ease factor based on the rating
//
//		EF':=EF+(0.1-(5-q)*(0.08+(5-q)*0.02))
//
// EF will never be lower than 1.3 or higher than 2.5
func UpdateEF(ef float64, rating int64) float64 {
	rc := float64(5 - rating)
	nef := ef + (0.1 - rc*(0.08+rc*0.02))
	if nef > MaxEF {
		nef = MaxEF
	} else if nef < MinEF {
		nef = MinEF
	}
	return nef
}

func UpdateInterval(days int64, ef float64) int64 {
	x := float64(days) * float64(ef)
	return int64(math.Ceil(x))
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

type CardReviewStore interface {
	AllByUserID(userID int64) ([]*CardReview, error)
	EnabledByUserID(userID int64) ([]*CardReview, error)
}

type ReviewStore interface {
	Insert(*Review) error
	ChangeEnabledForUserDeck(enabled bool, userID, deckID int64) error
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
