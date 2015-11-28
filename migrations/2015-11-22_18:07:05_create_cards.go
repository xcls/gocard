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
			CREATE TABLE cards (
				id serial PRIMARY KEY,
				front text NOT NULL CHECK(length(front) < 4000),
				back text NOT NULL CHECK(length(back) < 4000),
				created_at timestamp with time zone DEFAULT(current_timestamp)
			)`)
			return err
		},
		Down: func(ctx interface{}) error {
			c := ctx.(*pg.Context)
			_, err := c.Tx.Exec("DROP TABLE cards")
			return err
		},
	}
	Migrations.Add(migration)
}
