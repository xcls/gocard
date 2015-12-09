package seed

import (
	"github.com/mcls/gocard/stores"
	"github.com/mcls/gocard/stores/common"
)

func Run() error {
	var err error

	// Users
	user := &common.User{Email: "maartencls@gmail.com"}
	user.SetPassword("password")
	err = stores.Store.Users.Insert(user)
	if err != nil {
		return err
	}

	// Decks & Cards
	deck := &common.Deck{Name: "Effective Go"}
	err = stores.Store.Decks.Insert(deck)
	if err != nil {
		return err
	}

	cs := []struct {
		context string
		front   string
		back    string
	}{
		{
			context: "Formatting",
			front:   "Which tool should you use for formatting?",
			back:    "gofmt",
		},
		{
			context: "Commentary",
			front:   "Why should a doc comment always start with the name?",
			back:    "So it can be searched with grep.",
		},
	}
	for _, c := range cs {
		card := &common.Card{
			Context: c.context,
			Front:   c.front,
			Back:    c.back,
			DeckID:  deck.ID,
		}
		err = stores.Store.Cards.Insert(card)
		if err != nil {
			return err
		}
	}

	return nil
}
