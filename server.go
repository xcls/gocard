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

	r.HandleFunc("/", indexHandler)

	r.HandleFunc("/decks/new", NewDeckHandler)
	r.HandleFunc("/decks/{id:[0-9]+}", ShowDeckHandler)
	r.HandleFunc("/decks/{id:[0-9]+}/cards/new", NewCardHandler)

	port := ":8080"
	log.Printf("Starting server on %q\n", port)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(port, nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	decks, err := stores.Store.Decks.All()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	renderer.HTML(w, http.StatusOK, "home", map[string]interface{}{
		"decks": decks,
	})
}

type DeckForm struct {
	Name string
}

func (f *DeckForm) ToRecord() *stores.DeckRecord {
	return &stores.DeckRecord{Name: f.Name}
}

func NewDeckHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		renderer.HTML(w, http.StatusOK, "decks/new", nil)
	} else {
		if err := r.ParseForm(); err != nil {
			http.Error(w, fmt.Sprintf("Can't parse form: %s", err.Error()), 500)
			return
		}
		deck := new(DeckForm)
		err := decoder.Decode(deck, r.PostForm)
		if err != nil {
			fmt.Fprintf(w, "%q \n", r.PostForm)
			http.Error(w, fmt.Sprintf("Can't decode: %s", err.Error()), 500)
			return
		}

		record := deck.ToRecord()
		if err := stores.Store.Decks.Insert(record); err != nil {
			http.Error(w, fmt.Sprintf("Can't persist: %s", err.Error()), 500)
			return
		}
		fmt.Fprintf(w, "%q \n", deck)
		fmt.Fprintf(w, "%q \n", record)
	}
}

func ShowDeckHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "id arg: "+err.Error(), 500)
		return
	}
	deck, err := stores.Store.Decks.Find(id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	cards, err := stores.Store.Cards.AllByDeckId(id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	renderer.HTML(w, http.StatusOK, "decks/show", map[string]interface{}{
		"Deck":  deck,
		"Cards": cards,
	})
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

func NewCardHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	deck, err := stores.Store.Decks.Find(id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if deck == nil {
		http.Error(w, fmt.Sprintf("Can't find deck with id %d", id), 500)
		return
	}

	card := new(CardForm)
	if r.Method == "GET" {
		renderer.HTML(w, http.StatusOK, "cards/new", map[string]interface{}{
			"Deck": deck,
			"Card": card,
		})
	} else {
		if err := r.ParseForm(); err != nil {
			http.Error(w, fmt.Sprintf("Can't parse form: %s", err.Error()), 500)
			return
		}
		err := decoder.Decode(card, r.PostForm)
		if err != nil {
			fmt.Fprintf(w, "%q \n", r.PostForm)
			http.Error(w, fmt.Sprintf("Can't decode: %s", err.Error()), 500)
			return
		}

		record := card.ToRecord()
		record.DeckId = deck.Id
		fmt.Println(card)
		fmt.Println(record)
		if err := stores.Store.Cards.Insert(record); err != nil {
			http.Error(w, fmt.Sprintf("Can't persist: %s", err.Error()), 500)
			return
		}
		http.Redirect(w, r,
			fmt.Sprintf("/decks/%d", record.DeckId),
			http.StatusFound)
	}
}
