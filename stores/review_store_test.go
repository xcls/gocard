package stores

import (
	"testing"

	"github.com/mcls/gocard/stores/common"
)

type reviewStoreFixture struct {
	user    *common.User
	deckOne *common.Deck
	deckTwo *common.Deck
	cards   []*common.Card
}

func setupReviewStoreFixture(t *testing.T, s *common.Store) *reviewStoreFixture {
	data := new(reviewStoreFixture)
	data.user = newUser("maartencls@gmail.com")
	if err := s.Users.Insert(data.user); err != nil {
		t.Fatal(err)
	}

	data.deckOne = newDeck("TestDeck")
	data.deckTwo = newDeck("Deck Two")
	for _, deck := range []*common.Deck{data.deckOne, data.deckTwo} {
		if err := s.Decks.Insert(deck); err != nil {
			t.Fatal(err)
		}
	}

	data.cards = []*common.Card{
		newCard(data.deckOne.ID, "Test 1"),
		newCard(data.deckOne.ID, "Test 2"),
		newCard(data.deckTwo.ID, "DISABLED"),
	}
	for _, card := range data.cards {
		if err := s.Cards.Insert(card); err != nil {
			t.Fatal(err)
		}
	}
	return data
}

func TestReviewStoreInsert(t *testing.T) {
	for _, s := range setupStores(t) {
		user := newUser("maartencls@gmail.com")
		if err := s.Users.Insert(user); err != nil {
			t.Fatal(err)
		}

		deck := newDeck("TestDeck")
		if err := s.Decks.Insert(deck); err != nil {
			t.Fatal(err)
		}

		card := newCard(deck.ID, "Test")
		if err := s.Cards.Insert(card); err != nil {
			t.Fatal(err)
		}

		review := &common.Review{
			Enabled:    true,
			EaseFactor: 1.4,
			Interval:   3,
			CardID:     card.ID,
			UserID:     user.ID,
		}
		if err := s.Reviews.Insert(review); err != nil {
			t.Fatal(err)
		}

		cardReviews, err := s.CardReviews.AllByUserID(user.ID)
		if err != nil {
			t.Fatal(err)
		}
		if len(cardReviews) != 1 {
			t.Fatalf("Only expected 1 record: %+v", cardReviews)
		}
		got := cardReviews[0]
		if got.CardID != review.CardID {
			t.Fatalf("CardID mismatch: %+v != %+v", got.CardID, review.CardID)
		}
		if got.DeckName != deck.Name {
			t.Fatalf("Deck name mismatch: %+v != %+v", got.DeckName, deck.Name)
		}
		if got.Enabled != true {
			t.Fatalf("Expected Enabled to be true, was %+v", got.Enabled)
		}
	}
}

func TestReviewStoreInsert_CantDuplicateCardIDPerUserID(t *testing.T) {
	for _, s := range setupStores(t) {
		user := newUser("maartencls@gmail.com")
		if err := s.Users.Insert(user); err != nil {
			t.Fatal(err)
		}

		deck := newDeck("TestDeck")
		if err := s.Decks.Insert(deck); err != nil {
			t.Fatal(err)
		}

		card := newCard(deck.ID, "Test")
		if err := s.Cards.Insert(card); err != nil {
			t.Fatal(err)
		}

		review := &common.Review{
			Enabled:    true,
			EaseFactor: 1.4,
			Interval:   3,
			CardID:     card.ID,
			UserID:     user.ID,
		}
		if err := s.Reviews.Insert(review); err != nil {
			t.Fatal(err)
		}
		if err := s.Reviews.Insert(review); err == nil {
			t.Fatal("UserID and CardID should be unique")
		}
	}
}

func TestEnableAllForDeckID(t *testing.T) {
	for _, s := range setupStores(t) {
		fx := setupReviewStoreFixture(t, store)
		err := s.Reviews.EnableAllForDeckID(fx.user.ID, fx.deckOne.ID)
		if err != nil {
			t.Fatal(err)
		}
		cardReviews, err := s.CardReviews.AllByUserID(fx.user.ID)
		if err != nil {
			t.Fatal(err)
		}
		if len(cardReviews) != 2 {
			t.Fatalf("Expected 2 records: %+q", cardReviews)
		}

		// Add another card after others have been enabled
		if err := s.Cards.Insert(newCard(fx.deckOne.ID, "Test 3")); err != nil {
			t.Fatal(err)
		}
		err = s.Reviews.EnableAllForDeckID(fx.user.ID, fx.deckOne.ID)
		if err != nil {
			t.Fatal(err)
		}
		cardReviews, err = s.CardReviews.AllByUserID(fx.user.ID)
		if err != nil {
			t.Fatal(err)
		}
		if len(cardReviews) != 3 {
			t.Fatalf("Expected 3 records: %+q", cardReviews)
		}
		for _, x := range cardReviews {
			if x.Enabled != true {
				t.Fatalf("Expected Enabled to be true: %+v", x)
			}
		}
	}
}

func TestDisableAllForDeckID(t *testing.T) {
	var err error
	for _, s := range setupStores(t) {
		fx := setupReviewStoreFixture(t, store)
		err = s.Reviews.EnableAllForDeckID(fx.user.ID, fx.deckOne.ID)
		if err != nil {
			t.Fatal(err)
		}
		err = s.Reviews.DisableAllForDeckID(fx.user.ID, fx.deckOne.ID)
		if err != nil {
			t.Fatal(err)
		}
		crs, err := s.CardReviews.AllByUserID(fx.user.ID)
		if err != nil {
			t.Fatal(err)
		}
		if len(crs) != 2 {
			t.Fatalf("Wrong length: %+v", crs)
		}
		for _, x := range crs {
			if x.Enabled != false {
				t.Fatalf("Expected Enabled to be false: %+v", x)
			}
		}
	}
}
