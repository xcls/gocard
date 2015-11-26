package stores

import (
	"testing"

	"github.com/mcls/gocard/dbutil"
)

func resetDatabase(t *testing.T) {
	db, err := dbutil.Connect()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec("DELETE FROM cards;"); err != nil {
		t.Fatal(err)
	}
}

func TestCardRecord_Insert(t *testing.T) {
	resetDatabase(t)
	card := &CardRecord{
		Front: "Hello",
		Back:  "World",
	}
	if err := Store().Cards.Insert(card); err != nil {
		t.Fatal(err)
	}
	cs, err := Store().Cards.All()
	if err != nil {
		t.Fatal(err)
	}
	lastCard := cs[0]
	if lastCard.Id != card.Id {
		t.Fatalf("Cards not equal. \n%d \n!= \n%d\n", card.Id, lastCard.Id)
	}
	if lastCard.CreatedAt.Unix() != card.CreatedAt.Unix() {
		t.Fatalf("Cards not equal. \n%q \n!= \n%q\n", card.CreatedAt, lastCard.CreatedAt)
	}
}
