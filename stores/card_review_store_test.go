package stores

import (
	"testing"
	"time"

	"github.com/mcls/gocard/stores/common"
)

func TestCardReviewStore_DueAt(t *testing.T) {
	for _, s := range setupStores(t) {
		user := newUser("maartencls@gmail.com")
		if err := s.Users.Insert(user); err != nil {
			t.Fatal(err)
		}

		deck := newDeck("Deck One")
		if err := s.Decks.Insert(deck); err != nil {
			t.Fatal(err)
		}

		cards := []*common.Card{
			newCard(deck.ID, "Test 1"),
			newCard(deck.ID, "Test 2"),
			newCard(deck.ID, "Test 3"),
		}
		for _, card := range cards {
			if err := s.Cards.Insert(card); err != nil {
				t.Fatal(err)
			}
		}

		dueCards, err := s.CardReviews.DueAt(user.ID, time.Now())
		if err != nil {
			t.Fatal(err)
		}
		if len(dueCards) != 0 {
			t.Fatal("Expected no due cards")
		}

		if err := s.Reviews.ChangeEnabledForUserDeck(true, user.ID, deck.ID); err != nil {
			t.Fatal(err)
		}
		timeChecks := []struct {
			ts   time.Time
			want int
		}{
			{time.Now(), 3},
			{time.Now().Add(-100 * time.Hour), 3},
		}
		for i, c := range timeChecks {
			dueCards, err = s.CardReviews.DueAt(user.ID, c.ts)
			if err != nil {
				t.Fatal(err)
			}
			if len(dueCards) != c.want {
				t.Fatalf("(check %d) Expected %d due cards, got %d",
					i, c.want, len(dueCards))
			}
		}

		answerChecks := []struct {
			ts         time.Time
			lastAnswer *common.Answer
			want       int
		}{
			{
				ts:         time.Now(),
				lastAnswer: &common.Answer{Rating: 1, CardID: cards[0].ID, UserID: user.ID},
				want:       3,
			},
			{
				ts:         time.Now().Add(-100 * time.Hour),
				lastAnswer: &common.Answer{Rating: 1, CardID: cards[0].ID, UserID: user.ID},
				want:       3,
			},
			{
				ts:         time.Now().Add(-100 * time.Hour),
				lastAnswer: &common.Answer{Rating: 2, CardID: cards[0].ID, UserID: user.ID},
				want:       3,
			},
			{
				ts:         time.Now().Add(-100 * time.Hour),
				lastAnswer: &common.Answer{Rating: 3, CardID: cards[0].ID, UserID: user.ID},
				want:       2,
			},
			{
				ts:         time.Now().Add(-100 * time.Hour),
				lastAnswer: &common.Answer{Rating: 5, CardID: cards[0].ID, UserID: user.ID},
				want:       2,
			},
		}
		for i, c := range answerChecks {
			t.Logf("Running answerCheck[%d]", i)
			if err := s.Answers.Insert(c.lastAnswer); err != nil {
				t.Fatal(err)
			}
			dueCards, err = s.CardReviews.DueAt(user.ID, c.ts)
			if err != nil {
				t.Fatal(err)
			}
			if len(dueCards) != c.want {
				t.Fatalf("Expected %d due cards, got %d", c.want, len(dueCards))
			}
		}
	}
}
