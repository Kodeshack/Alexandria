package routes

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"alexandria.app/models"
	"alexandria.app/view"
)

func LoginRoutes(r *mux.Router, config *models.Config, userStorage models.UserStorage, sessionStorage *models.SessionStorage) {
	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		v := view.New("layout", "login", config)

		if err := v.Render(w, nil, nil); err != nil {
			log.Print(err)
			view.RenderErrorView("Failed to render login view.", http.StatusInternalServerError, config, nil, w)
		}
	}).Methods(http.MethodGet)

	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		email := strings.TrimSpace(r.FormValue("email"))
		password := r.FormValue("password")

		if len(email) == 0 || len(password) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			v := view.New("layout", "login", config)

			if err := v.Render(w, nil, nil); err != nil {
				log.Print(err)
				view.RenderErrorView("Failed to render login view.", http.StatusInternalServerError, config, nil, w)
				return
			}

			return
		}

		user := userStorage.CheckUserLogin(email, password)
		if user == nil {
			w.WriteHeader(http.StatusUnauthorized)
			v := view.New("layout", "login", config)

			if err := v.Render(w, nil, nil); err != nil {
				log.Print(err)
				view.RenderErrorView("Failed to render login view.", http.StatusInternalServerError, config, nil, w)
				return
			}
			return
		}

		session := models.NewSession(user)

		sessionStorage.AddSession(session)

		http.SetCookie(w, session.Cookie(false))

		http.Redirect(w, r, "/", http.StatusFound)
	}).Methods(http.MethodPost)
}
