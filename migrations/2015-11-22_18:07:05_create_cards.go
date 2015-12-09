package migrations

import (
	"github.com/mcls/nomad"
	"github.com/mcls/nomad/pg"
)

func init() {
	migration := &nomad.Migration{
		Version: "2015-11-22_18:07:05",
		Up: func(ctx interface{}) error {
			c := ctx.(*pg.Context)
			_, err := c.Tx.Exec(`
			CREATE TABLE decks (
				id serial PRIMARY KEY,
				name text NOT NULL CHECK(length(name) < 300),
				created_at timestamp with time zone DEFAULT(current_timestamp)
			)`)
			if err != nil {
				return err
			}
			_, err = c.Tx.Exec(`
			CREATE TABLE cards (
				id serial PRIMARY KEY,
				context text NOT NULL CHECK(length(context) < 300),
				front text NOT NULL CHECK(length(front) < 4000),
				back text NOT NULL CHECK(length(back) < 4000),
				deck_id integer REFERENCES decks ON DELETE CASCADE,
				created_at timestamp with time zone DEFAULT(current_timestamp)
			)`)
			return err
		},
		Down: func(ctx interface{}) error {
			c := ctx.(*pg.Context)
			if _, err := c.Tx.Exec("DROP TABLE IF EXISTS cards"); err != nil {
				return err
			}
			_, err := c.Tx.Exec("DROP TABLE IF EXISTS decks")
			return err
		},
	}
	Migrations.Add(migration)
}
