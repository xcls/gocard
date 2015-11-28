package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/Schema"
	"github.com/mcls/gocard/stores"
	"github.com/unrolled/render"
)

var decoder = schema.NewDecoder()

var renderer = render.New(render.Options{
	Directory: "templates",
	Layout:    "layout",
})

func startServer() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/cards/new", NewCardHandler)
	http.HandleFunc("/decks/new", NewDeckHandler)

	port := ":8080"
	log.Printf("Starting server on %q\n", port)
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
		"lol":   "lolsies",
		"hah":   []string{"hello", "hello"},
	})
}

func NewCardHandler(w http.ResponseWriter, r *http.Request) {
	renderer.HTML(w, http.StatusOK, "cards/new", nil)
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
