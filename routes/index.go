package routes

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"alexandria.app/view"
)

func IndexRoutes(r *mux.Router, templateDir string) {
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		v := view.New("layout", "index", templateDir, nil)

		if err := v.Render(w); err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}).Methods("GET")
}
