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
			view.RenderErrorView("Failed to render setup view.", http.StatusInternalServerError, config, nil, w)
			return
		}
	}).Methods(http.MethodGet)

	r.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		email := r.FormValue("email")
		displayName := strings.TrimSpace(r.FormValue("display_name"))
		password := r.FormValue("password")
		passwordConfirmation := r.FormValue("confirm_password")

		if len(email) == 0 || len(displayName) == 0 || len(password) == 0 || len(passwordConfirmation) == 0 {
			view.RenderErrorView("Email, display name, password or password confirmation empty.", http.StatusBadRequest, config, nil, w)
			return
		}

		if password != passwordConfirmation {
			view.RenderErrorView("Passwords don't match.", http.StatusBadRequest, config, nil, w)
			return
		}

		parsedEmail, err := mail.ParseAddress(email)
		if err != nil {
			view.RenderErrorView("Invalid email address.", http.StatusBadRequest, config, nil, w)
			return
		}
		email = parsedEmail.Address

		user, err := models.NewUser(email, displayName, password, true) // First user has to be an admin.
		if err != nil {
			log.Fatal(err)
			view.RenderErrorView("Failed to create new user.", http.StatusInternalServerError, config, nil, w)
			return
		}

		err = userStorage.AddUser(user)
		if err != nil {
			// User already exists or, very unlikely, a UUID collision.
			view.RenderErrorView("Failed to add user to user database.", http.StatusBadRequest, config, nil, w)
			return
		}

		err = userStorage.Save()
		if err != nil {
			log.Fatal(err)
			view.RenderErrorView("Failed to save user database.", http.StatusInternalServerError, config, nil, w)
			return
		}

		session := models.NewSession(user)
		sessionStorage.AddSession(session)
		http.SetCookie(w, session.Cookie(false))
		http.Redirect(w, r, "/", http.StatusFound)
	}).Methods(http.MethodPost)
}
