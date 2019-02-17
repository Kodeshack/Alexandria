package routes

import (
	"log"
	"net/http"
	"net/mail"
	"strings"

	"alexandria.app/models"
	"alexandria.app/view"
	"github.com/gorilla/mux"
)

func SetupRoutes(r *mux.Router, config *models.Config, userStorage models.UserStorage, sessionStorage *models.SessionStorage) {
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !userStorage.IsEmpty() {
				http.Redirect(w, r, "/", http.StatusFound)
			} else {
				next.ServeHTTP(w, r)
			}
		})
	})

	r.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		v := view.New("layout", "setup", config)
		if err := v.Render(w, nil, nil); err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}).Methods(http.MethodGet)

	r.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		email := r.FormValue("email")
		displayName := strings.TrimSpace(r.FormValue("display_name"))
		password := r.FormValue("password")
		passwordConfirmation := r.FormValue("confirm_password")

		if len(email) == 0 || len(displayName) == 0 || len(password) == 0 || len(passwordConfirmation) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if password != passwordConfirmation {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		parsedEmail, err := mail.ParseAddress(email)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		email = parsedEmail.Address

		user, err := models.NewUser(email, displayName, password, true) // First user has to be an admin.
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = userStorage.AddUser(user)
		if err != nil {
			// User already exists or, very unlikely, a UUID collision.
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = userStorage.Save()
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		session := models.NewSession(user)
		sessionStorage.AddSession(session)
		http.SetCookie(w, session.Cookie(false))
		http.Redirect(w, r, "/", http.StatusFound)
	}).Methods(http.MethodPost)
}
