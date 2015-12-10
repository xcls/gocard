package psql

import (
	"time"

	"github.com/mcls/gocard/stores/common"
)

const (
	SqlStart = `SELECT
	  r.id AS review_id,
	  r.enabled AS enabled,
	  r.ease_factor AS ease_factor,
	  r.interval AS interval,
	  r.due_on AS due_on,
	  r.user_id AS user_id,

	  c.id AS card_id,
	  c.context AS card_context,
	  c.front AS card_front,
	  c.back AS card_back,

	  d.id AS deck_id,
	  d.name AS deck_name
	FROM reviews r
	JOIN cards c ON c.id = r.card_id
	JOIN decks d ON d.id = c.deck_id`
)

type CardReviews dbmapStore

type CardReviewRecord struct {
	ID         int64     `db:"review_id"`
	Enabled    bool      `db:"enabled"`
	EaseFactor float64   `db:"ease_factor"`
	Interval   int64     `db:"interval"`
	DueOn      time.Time `db:"due_on"`
	UserID     int64     `db:"user_id"`

	CardID      int64  `db:"card_id"`
	CardContext string `db:"card_context"`
	CardFront   string `db:"card_front"`
	CardBack    string `db:"card_back"`

	DeckID   int64  `db:"deck_id"`
	DeckName string `db:"deck_name"`
}

func (r *CardReviewRecord) ToModel() *common.CardReview {
	return &common.CardReview{
		ID:         r.ID,
		Enabled:    r.Enabled,
		EaseFactor: r.EaseFactor,
		Interval:   r.Interval,
		DueOn:      r.DueOn,
		UserID:     r.UserID,

		CardID:      r.CardID,
		CardContext: r.CardContext,
		CardFront:   r.CardFront,
		CardBack:    r.CardBack,

		DeckID:   r.DeckID,
		DeckName: r.DeckName,
	}
}

func (r *CardReviewRecord) FromModel(m *common.CardReview) *CardReviewRecord {
	return &CardReviewRecord{
		ID:         m.ID,
		Enabled:    m.Enabled,
		EaseFactor: m.EaseFactor,
		Interval:   m.Interval,
		DueOn:      m.DueOn,
		UserID:     m.UserID,

		CardID:      m.CardID,
		CardContext: m.CardContext,
		CardFront:   m.CardFront,
		CardBack:    m.CardBack,

		DeckID:   m.DeckID,
		DeckName: m.DeckName,
	}
}

func (s *CardReviews) AllByUserID(userID int64) ([]*common.CardReview, error) {
	var rows []*CardReviewRecord
	sql := SqlStart + ` WHERE user_id = $1`
	_, err := s.DbMap.Select(&rows, sql, userID)
	if err != nil {
		return nil, err
	}
	return s.recordsToModels(rows), nil
}

func (s *CardReviews) EnabledByUserID(userID int64) ([]*common.CardReview, error) {
	var rows []*CardReviewRecord
	sql := SqlStart + ` WHERE user_id = $1 AND enabled = true`
	_, err := s.DbMap.Select(&rows, sql, userID)
	if err != nil {
		return nil, err
	}
	return s.recordsToModels(rows), nil
}

func (s *CardReviews) recordsToModels(rows []*CardReviewRecord) []*common.CardReview {
	models := make([]*common.CardReview, len(rows))
	for i, record := range rows {
		models[i] = record.ToModel()
	}
	return models
}
