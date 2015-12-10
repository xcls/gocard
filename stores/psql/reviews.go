package psql

import (
	"time"

	"github.com/mcls/gocard/stores/common"
	"gopkg.in/gorp.v1"
)

type Reviews dbmapStore

type ReviewRecord struct {
	ID         int64     `db:"id"`
	Enabled    bool      `db:"enabled"`
	EaseFactor float64   `db:"ease_factor"`
	Interval   int64     `db:"interval"`
	DueOn      time.Time `db:"due_on"`
	CardID     int64     `db:"card_id"`
	UserID     int64     `db:"user_id"`
	CreatedAt  time.Time `db:"created_at"`
}

func (r *ReviewRecord) PreInsert(s gorp.SqlExecutor) error {
	r.CreatedAt = time.Now()
	return nil
}

func (r *ReviewRecord) ToModel() *common.Review {
	return &common.Review{
		ID:         r.ID,
		Enabled:    r.Enabled,
		EaseFactor: r.EaseFactor,
		Interval:   r.Interval,
		DueOn:      r.DueOn,
		CardID:     r.CardID,
		UserID:     r.UserID,
		CreatedAt:  r.CreatedAt,
	}
}

func (r *ReviewRecord) FromModel(m *common.Review) *ReviewRecord {
	return &ReviewRecord{
		ID:         m.ID,
		Enabled:    m.Enabled,
		EaseFactor: m.EaseFactor,
		Interval:   m.Interval,
		DueOn:      m.DueOn,
		CardID:     m.CardID,
		UserID:     m.UserID,
		CreatedAt:  m.CreatedAt,
	}
}

func (s *Reviews) Insert(model *common.Review) error {
	record := new(ReviewRecord).FromModel(model)
	err := s.DbMap.Insert(record)
	if err != nil {
		return err
	}
	// Update original model with values from db
	*model = *record.ToModel()
	return nil
}
