package stores

import (
	"log"
	"time"

	"github.com/coopernurse/gorp"
	"github.com/mcls/gocard/dbutil"
)

type Stores struct {
	Cards *Cards
}

var stores = &Stores{
	Cards: &Cards{},
}

func Store() *Stores {
	return stores
}

// DB-backed store for Card
type Cards struct{}

type CardRecord struct {
	Id        int64     `db:"id"`
	Front     string    `db:"front"`
	Back      string    `db:"back"`
	CreatedAt time.Time `db:"created_at"`
	// UpdatedAt time.Time `db:"updated_at"`
}

// implement the PreInsert and PreUpdate hooks
func (c *CardRecord) PreInsert(s gorp.SqlExecutor) error {
	c.CreatedAt = time.Now()
	return nil
}

func (s *Cards) Insert(card *CardRecord) error {
	dbmap := initDb()
	return dbmap.Insert(card)
}

func (s *Cards) All() ([]*CardRecord, error) {
	var cards []*CardRecord
	dbmap := initDb()
	_, err := dbmap.Select(&cards, "SELECT * FROM cards ORDER BY created_at DESC")
	return cards, err
}

func initDb() *gorp.DbMap {
	db, err := dbutil.Connect()
	if err != nil {
		log.Fatal(err)
	}

	// construct a gorp DbMap
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}
	dbmap.AddTableWithName(CardRecord{}, "cards").SetKeys(true, "id")

	return dbmap
}
