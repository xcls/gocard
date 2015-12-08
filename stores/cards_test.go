package stores

import (
	"testing"

	"github.com/mcls/gocard/config"
	"github.com/mcls/gocard/dbutil"
	"github.com/mcls/gocard/migrations"
	nomadpg "github.com/mcls/nomad/pg"
	"github.com/stretchr/testify/assert"
)

func resetDatabase(t *testing.T) {
	db, err := dbutil.Connect(config.DatabaseTestURL())
	if err != nil {
		t.Fatal(err)
	}
	runner := nomadpg.NewRunner(db, migrations.Migrations)
	if err := runner.Run(); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec("DELETE FROM cards;"); err != nil {
		t.Fatal(err)
	}
}

func createDeck(t *testing.T, name string) *DeckRecord {
	deck := &DeckRecord{Name: name}
	if err := Store.Decks.Insert(deck); err != nil {
		t.Fatal(err)
	}
	return deck
}

func TestCards_Insert(t *testing.T) {
	resetDatabase(t)
	deck := createDeck(t, "Coding Knowledge")
	card := &Card{
		Context: "Programming",
		Front:   "Hello [...]",
		Back:    "Hello World",
		DeckID:  deck.ID,
	}
	assert.NoError(t, Store.Cards.Insert(card))
	cs, err := Store.Cards.All()
	if err != nil {
		t.Fatal(err)
	}
	lastCard := cs[0]
	assert.Equal(t, lastCard.ID, card.ID, "Card IDs don't match")
	assert.Equal(t, lastCard.CreatedAt.Unix(), card.CreatedAt.Unix(),
		"Card CreatedAt not equal")
}

func TestCards_Find(t *testing.T) {
	resetDatabase(t)
	deck := createDeck(t, "Coding Knowledge")
	card := &Card{
		Context: "Programming 2",
		Front:   "Hello [...]",
		Back:    "Hello World",
		DeckID:  deck.ID,
	}
	err := Store.Cards.Insert(card)
	assert.NoError(t, err)

	actual, err := Store.Cards.Find(card.ID)
	assert.NoError(t, err)
	assert.Equal(t, actual.ID, card.ID,
		"Card IDs not equal")
	assert.Equal(t, actual.CreatedAt.Unix(), card.CreatedAt.Unix(),
		"Card CreatedAt not equal")
}

func TestCards_Update(t *testing.T) {
	resetDatabase(t)
	var err error
	deck := createDeck(t, "Coding Knowledge")
	card := &Card{
		Context: "Programming 2",
		Front:   "Hello [...]",
		Back:    "Hello World",
		DeckID:  deck.ID,
	}
	err = Store.Cards.Insert(card)
	assert.NoError(t, err)

	card.Context = "Hacking 101"
	err = Store.Cards.Update(card)
	assert.NoError(t, err)

	actual, err := Store.Cards.Find(card.ID)
	assert.Equal(t, actual.Context, "Hacking 101")
}
