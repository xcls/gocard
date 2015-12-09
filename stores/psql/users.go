package psql

import (
	"time"

	"github.com/mcls/gocard/stores/common"
)

type Users dbmapStore

type UserRecord struct {
	ID                int64     `db:"id"`
	Email             string    `db:"email"`
	EncryptedPassword string    `db:"encrypted_password"`
	CreatedAt         time.Time `db:"created_at"`
	UpdatedAt         time.Time `db:"updated_at"`
}

func (r *UserRecord) ToModel() *common.User {
	return &common.User{
		ID:                r.ID,
		Email:             r.Email,
		EncryptedPassword: r.EncryptedPassword,
		CreatedAt:         r.CreatedAt,
		UpdatedAt:         r.UpdatedAt,
	}
}

func (r *UserRecord) FromModel(m *common.User) *UserRecord {
	return &UserRecord{
		ID:                m.ID,
		Email:             m.Email,
		EncryptedPassword: m.EncryptedPassword,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
	}
}

func (s *Users) Insert(model *common.User) error {
	record := new(UserRecord).FromModel(model)
	err := s.DbMap.Insert(record)
	if err != nil {
		return err
	}
	// Update original model with values from db
	*model = *record.ToModel()
	return nil
}

func (s *Users) Find(id int64) (*common.User, error) {
	var record *UserRecord
	err := s.DbMap.SelectOne(&record, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return record.ToModel(), nil
}

func (s *Users) Authenticate(email, password string) (*common.User, error) {
	var record *UserRecord
	err := s.DbMap.SelectOne(
		&record, "SELECT * FROM users WHERE email = $1", email,
	)
	if err != nil {
		return nil, err
	}
	model := record.ToModel()
	err = model.ComparePassword(password)
	if err != nil {
		return nil, err
	}
	return model, nil
}
