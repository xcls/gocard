package stores

import (
	"log"
	"time"

	"github.com/go-gorp/gorp"
	"github.com/mcls/gocard/config"
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
	ID        int64     `db:"id"`
	Context   string    `db:"context"`
	Front     string    `db:"front"`
	Back      string    `db:"back"`
	DeckID    int64     `db:"deck_id"`
	CreatedAt time.Time `db:"created_at"`
	// UpdatedAt time.Time `db:"updated_at"`
}

func (r *CardRecord) ToModel() *Card {
	return &Card{
		ID:        r.ID,
		Context:   r.Context,
		Front:     r.Front,
		Back:      r.Back,
		DeckID:    r.DeckID,
		CreatedAt: r.CreatedAt,
	}
}

func (r *CardRecord) FromModel(m *Card) *CardRecord {
	return &CardRecord{
		ID:        m.ID,
		Context:   m.Context,
		Front:     m.Front,
		Back:      m.Back,
		DeckID:    m.DeckID,
		CreatedAt: m.CreatedAt,
	}
}

type Card struct {
	ID        int64
	Context   string
	Front     string
	Back      string
	DeckID    int64
	CreatedAt time.Time
}

// implement the PreInsert and PreUpdate hooks
func (c *CardRecord) PreInsert(s gorp.SqlExecutor) error {
	c.CreatedAt = time.Now()
	return nil
}

type DeckRecord struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	// UpdatedAt time.Time `db:"updated_at"`
}

func (r *DeckRecord) PreInsert(s gorp.SqlExecutor) error {
	r.CreatedAt = time.Now()
	return nil
}

func (s *Cards) Insert(model *Card) error {
	record := new(CardRecord).FromModel(model)
	err := dbmap.Insert(record)
	if err != nil {
		return err
	}
	// Update original model with values from db
	*model = *record.ToModel()
	return nil
}

func (s *Cards) Update(model *Card) error {
	record := new(CardRecord).FromModel(model)
	_, err := dbmap.Update(record)
	if err != nil {
		return err
	}
	// Update original model with values from db
	*model = *record.ToModel()
	return nil
}

func (s *Cards) All() ([]*CardRecord, error) {
	var cards []*CardRecord
	_, err := dbmap.Select(&cards, "SELECT * FROM cards ORDER BY created_at DESC")
	return cards, err
}

func (s *Cards) AllByDeckID(id int) ([]*CardRecord, error) {
	var cards []*CardRecord
	_, err := dbmap.Select(
		&cards,
		"SELECT * FROM cards WHERE deck_id = $1 ORDER BY created_at DESC",
		id,
	)
	return cards, err
}

func (s *Cards) Find(id int64) (*Card, error) {
	var record *CardRecord
	err := dbmap.SelectOne(&record, "SELECT * FROM cards WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return record.ToModel(), nil
}

func (s *Decks) Insert(deck *DeckRecord) error {
	return dbmap.Insert(deck)
}

func (s *Decks) Find(id int64) (*DeckRecord, error) {
	var deck *DeckRecord
	err := dbmap.SelectOne(&deck, "SELECT * FROM decks WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return deck, nil
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
	dbmap.TraceOn("(SQL)", config.DefaultLogger())
	dbmap.AddTableWithName(CardRecord{}, "cards").SetKeys(true, "id")
	dbmap.AddTableWithName(DeckRecord{}, "decks").SetKeys(true, "id")

	return dbmap
}

func init() {
	dbmap = initDb()
}
