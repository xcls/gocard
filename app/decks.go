package app

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/mcls/gocard/stores"
	"github.com/mcls/gocard/stores/common"
	"github.com/mcls/gocard/valid"
)

type DeckForm struct {
	Name string
}

func (f *DeckForm) ToRecord() *common.Deck {
	return &common.Deck{Name: f.Name}
}

func (f *DeckForm) Validate() []error {
	vd := valid.NewValidator()
	vd.ValidateMinLength("Name", f.Name, 1)
	return vd.Errors()
}

func NewDeckHandler(rc *RequestContext) error {
	if rc.Request.Method == "GET" {
		return rc.HTML(http.StatusOK, "decks/new", tplVars{
			"DeckForm":   new(DeckForm),
			"DeckErrors": []error{},
		})
	}

	deck := new(DeckForm)
	if err := decodeForm(deck, rc.Request); err != nil {
		return err
	}

	if errs := deck.Validate(); len(errs) != 0 {
		return rc.HTML(http.StatusOK, "decks/new", tplVars{
			"DeckForm":   deck,
			"DeckErrors": errs,
		})
	}

	record := deck.ToRecord()
	if err := stores.Store.Decks.Insert(record); err != nil {
		return err
	}
	if err := rc.AddFlash("Saved Deck: " + record.Name); err != nil {
		return err
	}
	http.Redirect(rc.Writer, rc.Request,
		fmt.Sprintf("/decks/%d", record.ID),
		http.StatusFound)
	return nil
}

func ShowDeckHandler(rc *RequestContext) error {
	id, err := strconv.Atoi(rc.Vars()["id"])
	if err != nil {
		return err
	}
	deck, err := stores.Store.Decks.Find(int64(id))
	if err != nil {
		return err
	}
	cards, err := stores.Store.Cards.AllByDeckID(id)
	if err != nil {
		return err
	}
	return rc.HTML(http.StatusOK, "decks/show", tplVars{
		"Deck":  deck,
		"Cards": cards,
	})
}

func ToggleDeckHandler(rc *RequestContext) error {
	id, err := strconv.Atoi(rc.Vars()["id"])
	if err != nil {
		return err
	}
	deck, err := stores.Store.Decks.Find(int64(id))
	if err != nil {
		return err
	}

	q := rc.Request.URL.Query()
	enabled := q.Get("disable") != "true"
	err = rc.Store.Reviews.ChangeEnabledForUserDeck(enabled, rc.CurrentUser.ID, deck.ID)
	if err != nil {
		return err
	}

	var msg string
	if enabled {
		msg = "Cards in deck enabled"
	} else {
		msg = "Cards in deck disabled"
	}
	url := fmt.Sprintf("/decks/%d", deck.ID)
	return rc.RedirectWithFlash(url, msg)
}
