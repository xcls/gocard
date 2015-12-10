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

		cardReviews, err := store.CardReviews.AllByUserID(user.ID)
		if err != nil {
			t.Fatal(err)
		}
		if len(cardReviews) != 1 {
			t.Fatalf("Only expeted 1 record: %+v", cardReviews)
		}
		got := cardReviews[0]
		if got.CardID != review.CardID {
			t.Fatalf("CardID mismatch: %+v != %+v", got.CardID, review.CardID)
		}
		if got.DeckName != deck.Name {
			t.Fatalf("Deck name mismatch: %+v != %+v", got.DeckName, deck.Name)
		}

	}
}
