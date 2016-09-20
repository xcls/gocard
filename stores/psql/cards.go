package psql

import (
	"time"

	"github.com/xcls/gocard/stores/common"
	"gopkg.in/gorp.v1"
)

// DB-backed store for Card
type Cards dbmapStore

type CardRecord struct {
	ID        int64     `db:"id"`
	Context   string    `db:"context"`
	Front     string    `db:"front"`
	Back      string    `db:"back"`
	DeckID    int64     `db:"deck_id"`
	CreatedAt time.Time `db:"created_at"`
	// UpdatedAt time.Time `db:"updated_at"`
}

// implement the PreInsert and PreUpdate hooks
func (c *CardRecord) PreInsert(s gorp.SqlExecutor) error {
	c.CreatedAt = time.Now()
	return nil
}

func (r *CardRecord) ToModel() *common.Card {
	return &common.Card{
		ID:        r.ID,
		Context:   r.Context,
		Front:     r.Front,
		Back:      r.Back,
		DeckID:    r.DeckID,
		CreatedAt: r.CreatedAt,
	}
}

func (r *CardRecord) FromModel(m *common.Card) *CardRecord {
	return &CardRecord{
		ID:        m.ID,
		Context:   m.Context,
		Front:     m.Front,
		Back:      m.Back,
		DeckID:    m.DeckID,
		CreatedAt: m.CreatedAt,
	}
}

func (s *Cards) Insert(model *common.Card) error {
	record := new(CardRecord).FromModel(model)
	err := s.DbMap.Insert(record)
	if err != nil {
		return err
	}
	// Update original model with values from db
	*model = *record.ToModel()
	return nil
}

func (s *Cards) Update(model *common.Card) error {
	record := new(CardRecord).FromModel(model)
	_, err := s.DbMap.Update(record)
	if err != nil {
		return err
	}
	// Update original model with values from db
	*model = *record.ToModel()
	return nil
}

func (s *Cards) All() ([]*common.Card, error) {
	var records []*CardRecord
	_, err := s.DbMap.Select(&records, "SELECT * FROM cards ORDER BY created_at DESC")
	return s.recordsToModels(records), err
}

func (s *Cards) recordsToModels(rows []*CardRecord) []*common.Card {
	models := make([]*common.Card, len(rows))
	for i, record := range rows {
		models[i] = record.ToModel()
	}
	return models
}

func (s *Cards) AllByDeckID(id int) ([]*common.Card, error) {
	var records []*CardRecord
	_, err := s.DbMap.Select(
		&records,
		"SELECT * FROM cards WHERE deck_id = $1 ORDER BY created_at DESC",
		id,
	)
	return s.recordsToModels(records), err
}

func (s *Cards) Find(id int64) (*common.Card, error) {
	var record *CardRecord
	err := s.DbMap.SelectOne(&record, "SELECT * FROM cards WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return record.ToModel(), nil
}
