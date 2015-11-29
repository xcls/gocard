package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/Schema"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/mcls/gocard/config"
	"github.com/mcls/gocard/stores"
	"github.com/unrolled/render"
)

var decoder = schema.NewDecoder()

var renderer = render.New(render.Options{
	Directory: "templates",
	Layout:    "layout",
})

var jar = sessions.NewCookieStore([]byte(config.CookieSecret()))

type tplVars map[string]interface{}

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
	renderHTML(w, r, http.StatusOK, "home", tplVars{
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
		return renderHTML(w, r, http.StatusOK, "decks/new", nil)
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
	return renderHTML(w, r, http.StatusOK, "decks/show", tplVars{
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
		return renderHTML(w, r, http.StatusOK, "cards/new", tplVars{
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
		if err := addFlash(w, r, "Saved Card: "+record.Context); err != nil {
			return err
		}
		http.Redirect(w, r,
			fmt.Sprintf("/decks/%d", record.DeckId),
			http.StatusFound)
	}
	return nil
}

func renderHTML(w http.ResponseWriter, r *http.Request, status int, tpl string, vars tplVars) error {
	if vars == nil {
		vars = tplVars{}
	}
	flashes, err := getFlashes(w, r)
	if err != nil {
		return err
	}
	vars["Flashes"] = flashes
	renderer.HTML(w, status, tpl, vars)
	return nil
}

func addFlash(w http.ResponseWriter, r *http.Request, msg string) error {
	session, err := jar.Get(r, "ses")
	if err != nil {
		return err
	}
	session.AddFlash(msg)
	if err := session.Save(r, w); err != nil {
		return err
	}
	return nil
}

func getFlashes(w http.ResponseWriter, r *http.Request) ([]interface{}, error) {
	session, err := jar.Get(r, "ses")
	if err != nil {
		return nil, err
	}
	flashes := session.Flashes()
	if err := session.Save(r, w); err != nil {
		return nil, err
	}
	return flashes, nil
}
