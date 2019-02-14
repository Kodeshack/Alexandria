package routes

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"alexandria.app/models"
	"alexandria.app/view"
)

func AdminRoutes(r *mux.Router, config *models.Config, userStorage models.UserStorage) {
	r.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		user := models.GetRequestUser(r)
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		if !user.Admin {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		v := view.New("layout", "admin", config)
		if err := v.Render(w, user, userStorage.GetUsers()); err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}).Methods(http.MethodGet)
}
