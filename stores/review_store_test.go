package stores

import (
	"testing"

	"github.com/mcls/gocard/stores/common"
)

func TestReviewStoreInsert(t *testing.T) {
	for _, store := range setupStores(t) {
		user := newUser("maartencls@gmail.com")
		if err := store.Users.Insert(user); err != nil {
			t.Fatal(err)
		}

		deck := newDeck("TestDeck")
		if err := store.Decks.Insert(deck); err != nil {
			t.Fatal(err)
		}

		card := newCard(deck.ID, "Test")
		if err := store.Cards.Insert(card); err != nil {
			t.Fatal(err)
		}

		review := &common.Review{
			EaseFactor: 1.4,
			Interval:   3,
			CardID:     card.ID,
			UserID:     user.ID,
		}
		if err := store.Reviews.Insert(review); err != nil {
			t.Fatal(err)
		}
		if review.ID <= 0 {
			t.Fatalf("Invalid review ID: %+v", review.ID)
		}
	}
}
