package routes

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"alexandria.app/models"
	"alexandria.app/view"
)

func AdminRoutes(r *mux.Router, config *models.Config, userStorage models.UserStorage, sessionStorage *models.SessionStorage) {
	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			v := view.New("layout", "login", config, nil)

			if err := v.Render(w); err != nil {
				log.Print(err)
				w.WriteHeader(http.StatusInternalServerError)
			}

			return
		}

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		email := strings.TrimSpace(r.FormValue("email"))
		password := r.FormValue("password")

		if len(email) == 0 || len(password) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			v := view.New("layout", "login", config, nil)

			if err := v.Render(w); err != nil {
				log.Print(err)
				return
			}

			return
		}

		user := userStorage.CheckUserLogin(email, password)
		if user == nil {
			w.WriteHeader(http.StatusUnauthorized)
			v := view.New("layout", "login", config, nil)

			if err := v.Render(w); err != nil {
				log.Print(err)
				return
			}
			return
		}

		session := models.NewSession(user)

		sessionStorage.AddSession(session)

		http.SetCookie(w, session.Cookie(false))

		http.Redirect(w, r, "/", http.StatusFound)
	})

	r.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		if session := models.GetRequestSession(r); session != nil {
			sessionStorage.RemoveSession(session)
			http.SetCookie(w, models.RemoveSessionCookie(false))
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}).Methods(http.MethodGet)

	r.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		session := models.GetRequestSession(r)
		if session == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		v := view.New("layout", "admin", config, userStorage.GetUsers())
		if err := v.Render(w); err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}).Methods(http.MethodGet)
}
