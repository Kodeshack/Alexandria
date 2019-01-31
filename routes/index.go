package routes

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"alexandria.app/models"
	"alexandria.app/view"
)

func IndexRoutes(r *mux.Router, templateDir string) {
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		type viewData struct {
			User *models.User
		}

		data := viewData{}

		if session := models.GetRequestSession(r); session != nil {
			data.User = session.User
		}

		v := view.New("layout", "index", templateDir, data)

		if err := v.Render(w); err != nil {
			log.Print(err)
			return
		}
	}).Methods("GET")
}
