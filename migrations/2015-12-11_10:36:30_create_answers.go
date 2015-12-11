package migrations

import (
	"github.com/mcls/nomad"
	"github.com/mcls/nomad/pg"
)

func init() {
	migration := &nomad.Migration{
		Version: "2015-12-11_10:36:30",
		Up: func(ctx interface{}) error {
			c := ctx.(*pg.Context)
			_, err := c.Tx.Exec(`
			CREATE TABLE answers (
				id serial PRIMARY KEY,
				rating smallint NOT NULL CHECK(rating >= 0 AND rating <= 5),
				card_id	integer	NOT NULL REFERENCES cards ON DELETE RESTRICT,
				user_id	integer	NOT NULL REFERENCES users ON DELETE CASCADE,
				created_at timestamp with time zone DEFAULT(current_timestamp)
			)`)
			return err
		},
		Down: func(ctx interface{}) error {
			c := ctx.(*pg.Context)
			_, err := c.Tx.Exec(`DROP TABLE answers`)
			return err
		},
	}
	Migrations.Add(migration)
}
