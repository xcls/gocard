package migrations

import (
	"database/sql"
	"log"

	"github.com/mcls/gocard/config"
	"github.com/mcls/nomad"
	"github.com/mcls/nomad/pg"
	// Setup postgres driver
	_ "github.com/lib/pq"
)

var Migrations *nomad.List

func init() {
	db, err := sql.Open("postgres", config.DatabaseUrl())
	if err != nil {
		log.Fatal(err)
	}
	Migrations = pg.NewList(db)
}
