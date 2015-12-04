package app

import (
	"log"
	"net/http"

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

func StartServer() {
	r := mux.NewRouter()

	r.HandleFunc("/", errorHandler(indexHandler))

	r.HandleFunc("/decks/new", errorHandler(NewDeckHandler))
	r.HandleFunc("/decks/{id:[0-9]+}", errorHandler(ShowDeckHandler))
	r.HandleFunc("/decks/{id:[0-9]+}/cards/new", errorHandler(NewCardHandler))
	r.HandleFunc("/cards/{id:[0-9]+}/edit", errorHandler(EditCardHandler))

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
