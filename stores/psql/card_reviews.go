package psql

import (
	"time"

	"github.com/xcls/gocard/stores/common"
	"gopkg.in/gorp.v1"
)

const (
	SqlBasic = "SELECT * FROM user_cards"
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

	LastAnswerRating int64         `db:"last_answer_rating"`
	LastAnswerAt     gorp.NullTime `db:"last_answer_at"`
}

func (r *CardReviewRecord) ToModel() *common.CardReview {
	m := &common.CardReview{
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

		LastAnswerRating: r.LastAnswerRating,
	}
	if r.LastAnswerAt.Valid {
		m.LastAnswerAt = r.LastAnswerAt.Time
	}
	return m
}

func (s *CardReviews) Find(reviewID int64) (*common.CardReview, error) {
	var record *CardReviewRecord
	err := s.DbMap.SelectOne(
		&record,
		"SELECT * FROM user_cards WHERE review_id = $1",
		reviewID,
	)
	if err != nil {
		return nil, err
	}
	return record.ToModel(), nil
}

func (s *CardReviews) AllByUserID(userID int64) ([]*common.CardReview, error) {
	var rows []*CardReviewRecord
	sql := SqlBasic + ` WHERE user_id = $1`
	_, err := s.DbMap.Select(&rows, sql, userID)
	if err != nil {
		return nil, err
	}
	return s.recordsToModels(rows), nil
}

func (s *CardReviews) EnabledByUserID(userID int64) ([]*common.CardReview, error) {
	var rows []*CardReviewRecord
	sql := SqlBasic + ` WHERE user_id = $1 AND enabled = true`
	_, err := s.DbMap.Select(&rows, sql, userID)
	if err != nil {
		return nil, err
	}
	return s.recordsToModels(rows), nil
}

// DueAt returns all cards that are due. A card is due if we're passed its due
// date, or if the last answer for that card was lower than 3.
func (s *CardReviews) DueAt(userID int64, ts time.Time) ([]*common.CardReview, error) {
	var rows []*CardReviewRecord
	sql := "SELECT * FROM user_cards " +
		" WHERE user_id = $1 AND enabled = true" +
		" AND (due_on <= $2 OR coalesce(last_answer_rating, 0) < 3)"
	_, err := s.DbMap.Select(&rows, sql, userID, ts)
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
