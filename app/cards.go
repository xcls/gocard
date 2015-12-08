package app

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mcls/gocard/stores"
	"github.com/mcls/gocard/valid"
)

type CardForm struct {
	ID      int64
	Context string
	Front   string
	Back    string
}

func (f *CardForm) ToModel() *stores.Card {
	return &stores.Card{
		ID:      f.ID,
		Context: f.Context,
		Front:   f.Front,
		Back:    f.Back,
	}
}

func (f *CardForm) FromModel(m *stores.Card) *CardForm {
	*f = *&CardForm{
		ID:      m.ID,
		Context: m.Context,
		Front:   m.Front,
		Back:    m.Back,
	}
	return f
}

func (f *CardForm) Validate() []error {
	vd := valid.NewValidator()
	vd.ValidateMinLength("Context", f.Context, 1)
	vd.ValidateMinLength("Front", f.Front, 2)
	vd.ValidateMinLength("Back", f.Back, 1)
	return vd.Errors()
}

func NewCardHandler(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return err
	}
	deck, err := stores.Store.Decks.Find(int64(id))
	if err != nil {
		return err
	}
	if deck == nil {
		return err
	}

	card := new(CardForm)
	if r.Method == "GET" {
		return renderHTML(w, r, http.StatusOK, "cards/new", tplVars{
			"Deck": deck,
			"Card": card,
		})
	} else {
		err := decodeForm(card, r)
		if err != nil {
			return err
		}

		if errs := card.Validate(); len(errs) != 0 {
			return renderHTML(w, r, http.StatusOK, "cards/new", tplVars{
				"Deck":       deck,
				"Card":       card,
				"CardErrors": errs,
			})
		}
		model := card.ToModel()
		model.DeckID = deck.ID
		applog.Printf("Creating card: %+v \n", model)
		if err := stores.Store.Cards.Insert(model); err != nil {
			return err
		}
		if err := addFlash(w, r, "Saved Card: "+model.Context); err != nil {
			return err
		}
		http.Redirect(w, r,
			fmt.Sprintf("/decks/%d", model.DeckID),
			http.StatusFound)
	}
	return nil
}

func EditCardHandler(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return err
	}
	card, err := stores.Store.Cards.Find(int64(id))
	if err != nil {
		return err
	}
	deck, err := stores.Store.Decks.Find(int64(card.DeckID))
	if err != nil {
		return err
	}

	form := new(CardForm)
	if r.Method == "GET" {
		form.FromModel(card)
		return renderHTML(w, r, http.StatusOK, "cards/edit", tplVars{
			"Deck": deck,
			"Card": form,
		})
	} else {
		form.ID = card.ID
		err := decodeForm(form, r)
		if err != nil {
			return err
		}

		if errs := form.Validate(); len(errs) != 0 {
			return renderHTML(w, r, http.StatusOK, "cards/edit", tplVars{
				"Deck":       deck,
				"Card":       form,
				"CardErrors": errs,
			})
		}
		card := form.ToModel()
		card.DeckID = deck.ID

		if err := stores.Store.Cards.Update(card); err != nil {
			return err
		}
		addFlash(w, r, "Updated Card")
		http.Redirect(w, r,
			fmt.Sprintf("/decks/%d", card.DeckID),
			http.StatusFound)
	}
	return nil
}
