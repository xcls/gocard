package psql

import (
	"time"

	"github.com/mcls/gocard/stores/common"
)

type UserSessions dbmapStore

type UserSessionRecord struct {
	ID        int64     `db:"id"`
	UID       string    `db:"uid"`
	UserID    int64     `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
}

func (r *UserSessionRecord) ToModel() *common.UserSession {
	return &common.UserSession{
		ID:        r.ID,
		UID:       r.UID,
		UserID:    r.UserID,
		CreatedAt: r.CreatedAt,
	}
}

func (r *UserSessionRecord) FromModel(m *common.UserSession) *UserSessionRecord {
	return &UserSessionRecord{
		ID:        m.ID,
		UID:       m.UID,
		UserID:    m.UserID,
		CreatedAt: m.CreatedAt,
	}
}

func (s *UserSessions) Insert(model *common.UserSession) error {
	record := new(UserSessionRecord).FromModel(model)
	err := s.DbMap.Insert(record)
	if err != nil {
		return err
	}
	// Update original model with values from db
	*model = *record.ToModel()
	return nil
}

func (s *UserSessions) Find(uid string) (*common.UserSession, error) {
	var record *UserSessionRecord
	err := s.DbMap.SelectOne(&record, "SELECT * FROM user_sessions WHERE uid = $1", uid)
	if err != nil {
		return nil, err
	}
	return record.ToModel(), nil
}
