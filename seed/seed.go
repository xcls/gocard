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
		{
			context: "Naming",
			front:   "Are underscores or mixedCaps idiomatic for package names?",
			back:    "No",
		},
		{
			context: "Naming",
			front:   "The Reader in the bufio package is called Reader instead of BufReader to avoid what? (one word)",
			back:    "stutter",
		},
		{
			context: "Initialization",
			front:   "When are the init() functions of package called?",
			back:    "init is called after all the variable declarations in the package have evaluated their initializers",
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
