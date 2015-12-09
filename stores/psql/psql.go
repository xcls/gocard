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
		Cards: &Cards{DbMap: dbmap},
		Decks: &Decks{DbMap: dbmap},
		Users: &Users{DbMap: dbmap},
	}
}

// newDbMap creates and configures new gorp.DbMap
func newDbMap(db *sql.DB) *gorp.DbMap {
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}
	dbmap.TraceOn("(SQL)", config.DefaultLogger())
	dbmap.AddTableWithName(CardRecord{}, "cards").SetKeys(true, "id")
	dbmap.AddTableWithName(DeckRecord{}, "decks").SetKeys(true, "id")
	dbmap.AddTableWithName(UserRecord{}, "users").SetKeys(true, "id")
	return dbmap
}
