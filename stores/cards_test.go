package stores

import (
	"testing"

	"github.com/mcls/gocard/stores/common"
	"github.com/mcls/gocard/stores/psql"
	"github.com/mcls/gocard/testutil"
	"github.com/stretchr/testify/assert"
)

var store *common.Store

func resetDatabase(t *testing.T) {
	db := testutil.ConnectDB(t)
	testutil.ResetDB(t, db)
	store = psql.NewStore(db)
}

func createDeck(t *testing.T, name string) *common.Deck {
	deck := &common.Deck{Name: name}
	if err := store.Decks.Insert(deck); err != nil {
		t.Fatal(err)
	}
	return deck
}

func TestCards_Insert(t *testing.T) {
	resetDatabase(t)
	deck := createDeck(t, "Coding Knowledge")
	card := &common.Card{
		Context: "Programming",
		Front:   "Hello [...]",
		Back:    "Hello World",
		DeckID:  deck.ID,
	}
	assert.NoError(t, store.Cards.Insert(card))
	cs, err := store.Cards.All()
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
	card := &common.Card{
		Context: "Programming 2",
		Front:   "Hello [...]",
		Back:    "Hello World",
		DeckID:  deck.ID,
	}
	err := store.Cards.Insert(card)
	assert.NoError(t, err)

	actual, err := store.Cards.Find(card.ID)
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
	card := &common.Card{
		Context: "Programming 2",
		Front:   "Hello [...]",
		Back:    "Hello World",
		DeckID:  deck.ID,
	}
	err = store.Cards.Insert(card)
	assert.NoError(t, err)

	card.Context = "Hacking 101"
	err = store.Cards.Update(card)
	assert.NoError(t, err)

	actual, err := store.Cards.Find(card.ID)
	assert.Equal(t, actual.Context, "Hacking 101")
}
