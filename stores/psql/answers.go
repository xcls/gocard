package psql

import (
	"time"

	"github.com/xcls/gocard/stores/common"
	"gopkg.in/gorp.v1"
)

type Answers dbmapStore

type AnswerRecord struct {
	ID        int64     `db:"id"`
	Rating    int64     `db:"rating"`
	CardID    int64     `db:"card_id"`
	UserID    int64     `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
}

func (r *AnswerRecord) PreInsert(s gorp.SqlExecutor) error {
	r.CreatedAt = time.Now()
	return nil
}

func (r *AnswerRecord) ToModel() *common.Answer {
	return &common.Answer{
		ID:        r.ID,
		Rating:    r.Rating,
		CardID:    r.CardID,
		UserID:    r.UserID,
		CreatedAt: r.CreatedAt,
	}
}

func (r *AnswerRecord) FromModel(m *common.Answer) *AnswerRecord {
	return &AnswerRecord{
		ID:        m.ID,
		Rating:    m.Rating,
		CardID:    m.CardID,
		UserID:    m.UserID,
		CreatedAt: m.CreatedAt,
	}
}

func (s *Answers) Insert(model *common.Answer) error {
	record := new(AnswerRecord).FromModel(model)
	err := s.DbMap.Insert(record)
	if err != nil {
		return err
	}
	// Update original model with values from db
	*model = *record.ToModel()
	return nil
}
