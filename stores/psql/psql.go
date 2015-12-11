package psql

import (
	"database/sql"

	"github.com/mcls/gocard/config"
	"github.com/mcls/gocard/stores/common"
	"gopkg.in/gorp.v1"
)

type dbmapStore struct {
	DbMap *gorp.DbMap
}

func NewStore(db *sql.DB) *common.Store {
	dbmap := newDbMap(db)
	return &common.Store{
		Answers:      &Answers{DbMap: dbmap},
		CardReviews:  &CardReviews{DbMap: dbmap},
		Cards:        &Cards{DbMap: dbmap},
		Decks:        &Decks{DbMap: dbmap},
		Reviews:      &Reviews{DbMap: dbmap},
		UserSessions: &UserSessions{DbMap: dbmap},
		Users:        &Users{DbMap: dbmap},
	}
}

// newDbMap creates and configures new gorp.DbMap
func newDbMap(db *sql.DB) *gorp.DbMap {
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}
	dbmap.TraceOn("(SQL)", config.DefaultLogger())
	dbmap.TraceOff() // FIXME: Disable for tests only?
	dbmap.AddTableWithName(AnswerRecord{}, "answers").SetKeys(true, "id")
	dbmap.AddTableWithName(CardRecord{}, "cards").SetKeys(true, "id")
	dbmap.AddTableWithName(DeckRecord{}, "decks").SetKeys(true, "id")
	dbmap.AddTableWithName(ReviewRecord{}, "reviews").SetKeys(true, "id")
	dbmap.AddTableWithName(UserRecord{}, "users").SetKeys(true, "id")
	dbmap.AddTableWithName(UserSessionRecord{}, "user_sessions").SetKeys(true, "id")
	return dbmap
}
