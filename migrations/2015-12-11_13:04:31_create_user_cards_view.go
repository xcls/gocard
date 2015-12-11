package migrations

import (
	"github.com/mcls/nomad"
	"github.com/mcls/nomad/pg"
)

func init() {
	migration := &nomad.Migration{
		Version: "2015-12-11_13:04:31",
		Up: func(ctx interface{}) error {
			c := ctx.(*pg.Context)
			_, err := c.Tx.Exec(`
	CREATE OR REPLACE VIEW user_cards AS

WITH last_answers AS (
    SELECT
      ans.*,
      MAX(ans.created_at) OVER (PARTITION BY ans.card_id, ans.user_id) AS last_answer_at
    FROM answers ans
)
SELECT
  r.id AS review_id,
  r.enabled AS enabled,
  r.ease_factor AS ease_factor,
  r.interval AS interval,
  r.due_on AS due_on,
  r.user_id AS user_id,
  c.id AS card_id,
  c.context AS card_context,
  c.front AS card_front,
  c.back AS card_back,
  d.id AS deck_id,
  d.name AS deck_name,

  COALESCE(a.rating, 0) AS last_answer_rating,
  a.created_at AS last_answer_at

  FROM reviews r
  JOIN cards c ON c.id = r.card_id
  JOIN decks d ON d.id = c.deck_id
  LEFT JOIN last_answers a
    ON r.card_id = a.card_id AND r.user_id = a.user_id AND a.created_at = a.last_answer_at`)
			return err
		},
		Down: func(ctx interface{}) error {
			c := ctx.(*pg.Context)
			_, err := c.Tx.Exec("DROP VIEW user_cards")
			return err
		},
	}
	Migrations.Add(migration)
}
