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
	var decksWithCards = []struct {
		deck  *common.Deck
		cards []*common.Card
	}{
		{
			&common.Deck{Name: "Effective Go"},
			[]*common.Card{
				{
					Context: "Formatting",
					Front:   "Which tool should you use for formatting?",
					Back:    "gofmt",
				},
				{
					Context: "Commentary",
					Front:   "Why should a doc comment always start with the name?",
					Back:    "So it can be searched with grep.",
				},
				{
					Context: "Naming",
					Front:   "Are underscores or mixedCaps idiomatic for package names?",
					Back:    "No",
				},
				{
					Context: "Naming",
					Front:   "The Reader in the bufio package is called Reader instead of BufReader to avoid what? (one word)",
					Back:    "stutter",
				},
				{
					Context: "Initialization",
					Front:   "When are the init() functions of package called?",
					Back:    "init is called after all the variable declarations in the package have evaluated their initializers",
				},
			},
		},

		{
			&common.Deck{Name: "Spaced Repitition"},
			[]*common.Card{
				{
					Context: "Definitions",
					Front:   "What is the spacing effect?",
					Back:    "The phenomenon whereby animals (including humans) more easily remember or learn items when they are studied via spaced presentation rather than cramming",
				},
				{
					Context: "Definitions",
					Front:   "What is the spaced presentation?",
					Back:    "A learning technique that uses increasing intervals of time between subsequent reviews of already learned material ",
				},
			},
		},
	}
	for _, o := range decksWithCards {
		err = stores.Store.Decks.Insert(o.deck)
		if err != nil {
			return err
		}
		for _, card := range o.cards {
			card.DeckID = o.deck.ID
			err = stores.Store.Cards.Insert(card)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
