package stores

import (
	"testing"

	"github.com/xcls/gocard/stores/common"
)

type answerStoreFixture struct {
	user    *common.User
	deckOne *common.Deck
	cards   []*common.Card
}

func setupAnswerStoreFixture(t *testing.T, s *common.Store) *reviewStoreFixture {
	data := new(reviewStoreFixture)
	data.user = newUser("maartencls@gmail.com")
	if err := s.Users.Insert(data.user); err != nil {
		t.Fatal(err)
	}

	data.deckOne = newDeck("Deck One")
	for _, deck := range []*common.Deck{data.deckOne} {
		if err := s.Decks.Insert(deck); err != nil {
			t.Fatal(err)
		}
	}

	data.cards = []*common.Card{
		newCard(data.deckOne.ID, "Test 1"),
	}
	for _, card := range data.cards {
		if err := s.Cards.Insert(card); err != nil {
			t.Fatal(err)
		}
	}
	return data
}

func TestAnswerStoreInsert(t *testing.T) {
	for _, s := range setupStores(t) {
		fx := setupAnswerStoreFixture(t, s)
		answer := &common.Answer{
			Rating: 3,
			CardID: fx.cards[0].ID,
			UserID: fx.user.ID,
		}
		if err := s.Answers.Insert(answer); err != nil {
			t.Fatal(err)
		}
	}
}
