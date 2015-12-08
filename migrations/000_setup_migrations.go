package migrations

import (
	"log"

	"github.com/mcls/gocard/dbutil"
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
	db, err := dbutil.Connect()
	if err != nil {
		log.Fatal(err)
	}
	return nomadpg.NewRunner(db, Migrations)
}
