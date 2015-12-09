package app

import (
	"net/http"

	"github.com/codegangsta/negroni"
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
var applog = config.DefaultLogger()

func StartServer() {
	// Routes
	r := mux.NewRouter()
	r.HandleFunc("/", withContext(indexHandler))
	r.HandleFunc("/login", withContext(LoginHandler))
	r.HandleFunc("/logout", withContext(LogoutHandler))
	r.HandleFunc("/register", withContext(RegisterHandler))
	r.HandleFunc("/decks/new", withContext(NewDeckHandler))
	r.HandleFunc("/decks/{id:[0-9]+}", withContext(ShowDeckHandler))
	r.HandleFunc("/decks/{id:[0-9]+}/cards/new", withContext(NewCardHandler))
	r.HandleFunc("/cards/{id:[0-9]+}/edit", withContext(EditCardHandler))

	// Middleware
	httpLogger := &negroni.Logger{applog}
	n := negroni.New(
		negroni.NewRecovery(),
		httpLogger,
		negroni.Wrap(r),
	)

	port := ":8080"
	applog.Printf("Starting server on %q\n", port)
	applog.Fatal(http.ListenAndServe(port, n))
}

func withContext(f func(*RequestContext) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rc := &RequestContext{
			Writer:  w,
			Request: r,
			Store:   stores.Store,
		}

		determineCurrentUser(rc)

		// Request handler
		if err := f(rc); err != nil {
			handlerInternalError(rc, err)
		}
	}
}

func determineCurrentUser(rc *RequestContext) {
	// User Session Management
	session, err := jar.Get(rc.Request, "uid")
	if err != nil {
		handlerInternalError(rc, err)
	}
	uid := session.Values["uid"]
	if val, ok := uid.(string); ok == true && val != "" {
		userSession, err := rc.Store.UserSessions.Find(val)
		if err != nil {
			applog.Printf("Failed to set current user: %v", err)
			return
		}
		user, err := rc.Store.Users.Find(userSession.UserID)
		if err != nil {
			applog.Printf("Failed to set current user: %v", err)
			return
		}
		rc.CurrentUser = user
	}
}

func handlerInternalError(rc *RequestContext, err error) {
	rc.RenderInternalServerErrorHTML(err)
	applog.Printf("Internal Server Error on %q \n", rc.Request.RequestURI)
	applog.Printf("Error: %v \n", err)
}

func indexHandler(rc *RequestContext) error {
	decks, err := stores.Store.Decks.All()
	if err != nil {
		return err
	}
	rc.HTML(http.StatusOK, "home", tplVars{"decks": decks})
	return nil
}

func errorHandler(f func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			applog.Printf("Internal Server Error on %q \n", r.RequestURI)
			applog.Printf("Error: %v \n", err)
		}
	}
}
