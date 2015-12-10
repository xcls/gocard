package migrations

import (
	"github.com/mcls/nomad"
	"github.com/mcls/nomad/pg"
)

func init() {
	migration := &nomad.Migration{
		Version: "2015-12-10_09:12:16",
		Up: func(ctx interface{}) error {
			c := ctx.(*pg.Context)
			_, err := c.Tx.Exec(`
			CREATE TABLE reviews (
				id serial PRIMARY KEY,
				enabled		boolean			NOT NULL DEFAULT(true),
				ease_factor numeric(6, 4)	NOT NULL CHECK(ease_factor >= 1.3),
				interval	integer			NOT NULL CHECK(interval >= 0),
				due_on		timestamp with time zone DEFAULT(current_timestamp),
				card_id		integer			NOT NULL REFERENCES cards ON DELETE RESTRICT,
				user_id		integer			NOT NULL REFERENCES users ON DELETE CASCADE,
				created_at timestamp with time zone DEFAULT(current_timestamp),
				UNIQUE(card_id, user_id)
			)`)
			return err
		},
		Down: func(ctx interface{}) error {
			c := ctx.(*pg.Context)
			_, err := c.Tx.Exec(`DROP TABLE reviews`)
			return err
		},
	}
	Migrations.Add(migration)
}
