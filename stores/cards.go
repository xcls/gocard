package stores

import (
	"log"
	"time"

	"github.com/coopernurse/gorp"
	"github.com/mcls/gocard/dbutil"
)

type Stores struct {
	Cards *Cards
	Decks *Decks
}

var Store = &Stores{
	Cards: &Cards{},
	Decks: &Decks{},
}

// DB-backed store for Card
type Cards struct{}
type Decks struct{}

type CardRecord struct {
	Id        int64     `db:"id"`
	Context   string    `db:"context"`
	Front     string    `db:"front"`
	Back      string    `db:"back"`
	DeckId    int64     `db:"deck_id"`
	CreatedAt time.Time `db:"created_at"`
	// UpdatedAt time.Time `db:"updated_at"`
}

// implement the PreInsert and PreUpdate hooks
func (c *CardRecord) PreInsert(s gorp.SqlExecutor) error {
	c.CreatedAt = time.Now()
	return nil
}

type DeckRecord struct {
	Id        int64     `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	// UpdatedAt time.Time `db:"updated_at"`
}

func (r *DeckRecord) PreInsert(s gorp.SqlExecutor) error {
	r.CreatedAt = time.Now()
	return nil
}

func (s *Cards) Insert(card *CardRecord) error {
	return dbmap.Insert(card)
}

func (s *Cards) All() ([]*CardRecord, error) {
	var cards []*CardRecord
	_, err := dbmap.Select(&cards, "SELECT * FROM cards ORDER BY created_at DESC")
	return cards, err
}

func (s *Decks) Insert(deck *DeckRecord) error {
	return dbmap.Insert(deck)
}

func (s *Decks) All() ([]*DeckRecord, error) {
	var decks []*DeckRecord
	_, err := dbmap.Select(&decks, "SELECT * FROM decks ORDER BY created_at DESC")
	return decks, err
}

var dbmap *gorp.DbMap

// Setup gorp.DbMap
func initDb() *gorp.DbMap {
	db, err := dbutil.Connect()
	if err != nil {
		log.Fatal(err)
	}

	// construct a gorp DbMap
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}
	dbmap.AddTableWithName(CardRecord{}, "cards").SetKeys(true, "id")
	dbmap.AddTableWithName(DeckRecord{}, "decks").SetKeys(true, "id")

	return dbmap
}

func init() {
	dbmap = initDb()
}
