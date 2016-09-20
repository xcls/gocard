package psql

import (
	"time"

	"github.com/xcls/gocard/stores/common"
	"gopkg.in/gorp.v1"
)

// DB-backed store for Decks
type Decks dbmapStore

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

func (r *DeckRecord) ToModel() *common.Deck {
	return &common.Deck{
		ID:        r.ID,
		Name:      r.Name,
		CreatedAt: r.CreatedAt,
	}
}

func (r *DeckRecord) FromModel(m *common.Deck) *DeckRecord {
	return &DeckRecord{
		ID:        m.ID,
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
	}
}

func (s *Decks) Insert(model *common.Deck) error {
	record := new(DeckRecord).FromModel(model)
	err := s.DbMap.Insert(record)
	if err != nil {
		return err
	}
	// Update original model with values from db
	*model = *record.ToModel()
	return nil
}

func (s *Decks) Find(id int64) (*common.Deck, error) {
	var deck *DeckRecord
	err := s.DbMap.SelectOne(&deck, "SELECT * FROM decks WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return deck.ToModel(), nil
}

func (s *Decks) All() ([]*common.Deck, error) {
	var rows []*DeckRecord
	_, err := s.DbMap.Select(&rows, "SELECT * FROM decks ORDER BY created_at DESC")

	decks := make([]*common.Deck, len(rows))
	for i, record := range rows {
		decks[i] = record.ToModel()
	}
	return decks, err
}
