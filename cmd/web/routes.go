package main

import ( 
	"log"
	"net/http"

	"github.com/CloudyKit/jet/v6"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (a *application) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	if a.debug {
		mux.Use(middleware.Logger)
	}
	mux.Use(middleware.Recoverer)
	mux.Use(a.LoadSession)

	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {

		a.session.Put(r.Context(), "test", "Divine")
		a.infoLog.Println(a.session.Get(r.Context(), "test"))  

		err := a.render(w, r, "index", nil)
		if err != nil {
			log.Fatalln(err)
		}
	})

	mux.Get("/com", func(w http.ResponseWriter, r *http.Request) {

		a.infoLog.Println(a.session)
		a.infoLog.Println(a.session.Exists(r.Context(), "test"))
		a.infoLog.Println(a.session.GetString(r.Context(), "test")) 

		vars := make(jet.VarMap)
		vars.Set("test", a.session.GetString(r.Context(), "test"))

		err := a.render(w, r, "index", vars)
		if err != nil {
			log.Fatalln(err)
		}
	})

	return mux
}
