package migrations

import (
	"database/sql"
	"log"

	"github.com/mcls/gocard/config"
	"github.com/mcls/nomad"
	nomadpg "github.com/mcls/nomad/pg"
	// Setup postgres driver
	_ "github.com/lib/pq"
)

var Migrations *nomad.List

func init() {
	Migrations = nomad.NewList()
}

func NewRunner() *nomad.Runner {
	db, err := sql.Open("postgres", config.DatabaseURL())
	if err != nil {
		log.Fatal(err)
	}
	return nomadpg.NewRunner(db, Migrations)
}
