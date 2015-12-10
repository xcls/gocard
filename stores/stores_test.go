package stores

import (
	"testing"

	"github.com/mcls/gocard/stores/common"
	"github.com/mcls/gocard/stores/psql"
	"github.com/mcls/gocard/testutil"
)

// Store shared test helpers here

func setupStores(t *testing.T) []*common.Store {
	db := testutil.ConnectDB(t)
	testutil.ResetDB(t, db)
	return []*common.Store{psql.NewStore(db)}
}

func newUser(email string, optPass ...string) *common.User {
	user := &common.User{Email: email, EncryptedPassword: "NOT_SET"}
	if len(optPass) > 0 {
		user.SetPassword(optPass[0])
	}
	return user
}

func newDeck(name string) *common.Deck {
	return &common.Deck{Name: name}
}

func newCard(deckID int64, context string) *common.Card {
	return &common.Card{
		Context: context,
		Front:   "Is this a test question?",
		Back:    "Yes. Yes, it is.",
		DeckID:  deckID,
	}
}
