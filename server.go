package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/Schema"
	"github.com/gorilla/mux"
	"github.com/mcls/gocard/stores"
	"github.com/unrolled/render"
)

var decoder = schema.NewDecoder()

var renderer = render.New(render.Options{
	Directory: "templates",
	Layout:    "layout",
})

func startServer() {
	r := mux.NewRouter()

	r.HandleFunc("/", errorHandler(indexHandler))

	r.HandleFunc("/decks/new", errorHandler(NewDeckHandler))
	r.HandleFunc("/decks/{id:[0-9]+}", errorHandler(ShowDeckHandler))
	r.HandleFunc("/decks/{id:[0-9]+}/cards/new", errorHandler(NewCardHandler))

	port := ":8080"
	log.Printf("Starting server on %q\n", port)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(port, nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) error {
	decks, err := stores.Store.Decks.All()
	if err != nil {
		return err
	}
	renderer.HTML(w, http.StatusOK, "home", map[string]interface{}{
		"decks": decks,
	})
	return nil
}

func errorHandler(f func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("handling %q: %v \n", r.RequestURI, err)
		}
	}
}

type DeckForm struct {
	Name string
}

func (f *DeckForm) ToRecord() *stores.DeckRecord {
	return &stores.DeckRecord{Name: f.Name}
}

func NewDeckHandler(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		renderer.HTML(w, http.StatusOK, "decks/new", nil)
		return nil
	}

	if err := r.ParseForm(); err != nil {
		return err
	}
	deck := new(DeckForm)
	err := decoder.Decode(deck, r.PostForm)
	if err != nil {
		return err
	}

	record := deck.ToRecord()
	if err := stores.Store.Decks.Insert(record); err != nil {
		return err
	}
	fmt.Fprintf(w, "%q \n", deck)
	fmt.Fprintf(w, "%q \n", record)
	return nil
}

func ShowDeckHandler(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return err
	}
	deck, err := stores.Store.Decks.Find(id)
	if err != nil {
		return err
	}
	cards, err := stores.Store.Cards.AllByDeckId(id)
	if err != nil {
		return err
	}
	renderer.HTML(w, http.StatusOK, "decks/show", map[string]interface{}{
		"Deck":  deck,
		"Cards": cards,
	})
	return nil
}

type CardForm struct {
	Context string
	Front   string
	Back    string
}

func (f *CardForm) ToRecord() *stores.CardRecord {
	fmt.Println(f.Context)
	return &stores.CardRecord{
		Context: f.Context,
		Front:   f.Front,
		Back:    f.Back,
	}
}

func NewCardHandler(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return err
	}
	deck, err := stores.Store.Decks.Find(id)
	if err != nil {
		return err
	}
	if deck == nil {
		return err
	}

	card := new(CardForm)
	if r.Method == "GET" {
		renderer.HTML(w, http.StatusOK, "cards/new", map[string]interface{}{
			"Deck": deck,
			"Card": card,
		})
	} else {
		if err := r.ParseForm(); err != nil {
			return err
		}
		err := decoder.Decode(card, r.PostForm)
		if err != nil {
			return err
		}

		record := card.ToRecord()
		record.DeckId = deck.Id
		fmt.Println(card)
		fmt.Println(record)
		if err := stores.Store.Cards.Insert(record); err != nil {
			return err
		}
		http.Redirect(w, r,
			fmt.Sprintf("/decks/%d", record.DeckId),
			http.StatusFound)
	}
	return nil
}
