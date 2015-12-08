package app

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mcls/gocard/stores"
	"github.com/mcls/gocard/valid"
)

type DeckForm struct {
	Name string
}

func (f *DeckForm) ToRecord() *stores.DeckRecord {
	return &stores.DeckRecord{Name: f.Name}
}

func (f *DeckForm) Validate() []error {
	vd := valid.NewValidator()
	vd.ValidateMinLength("Name", f.Name, 1)
	return vd.Errors()
}

func NewDeckHandler(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return renderHTML(w, r, http.StatusOK, "decks/new", tplVars{
			"DeckForm":   new(DeckForm),
			"DeckErrors": []error{},
		})
	}

	deck := new(DeckForm)
	if err := decodeForm(deck, r); err != nil {
		return err
	}

	if errs := deck.Validate(); len(errs) != 0 {
		return renderHTML(w, r, http.StatusOK, "decks/new", tplVars{
			"DeckForm":   deck,
			"DeckErrors": errs,
		})
	}

	record := deck.ToRecord()
	if err := stores.Store.Decks.Insert(record); err != nil {
		return err
	}
	if err := addFlash(w, r, "Saved Deck: "+record.Name); err != nil {
		return err
	}
	http.Redirect(w, r,
		fmt.Sprintf("/decks/%d", record.ID),
		http.StatusFound)
	return nil
}

func ShowDeckHandler(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
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
	return renderHTML(w, r, http.StatusOK, "decks/show", tplVars{
		"Deck":  deck,
		"Cards": cards,
	})
}
