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

func UserRoutes(r *mux.Router, config *models.Config, userStorage models.UserStorage, sessionStorage *models.SessionStorage) {
	r.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		session := models.GetRequestSession(r)

		if session == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		v := view.New("layout", "user", config, session.User)

		if err := v.Render(w); err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}).Methods(http.MethodGet)

	r.HandleFunc("/user/update", func(w http.ResponseWriter, r *http.Request) {
		session := models.GetRequestSession(r)
		user := session.User

		if session == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		email := r.FormValue("email")
		displayName := strings.TrimSpace(r.FormValue("display_name"))

		if len(email) == 0 || len(displayName) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		parsedEmail, err := mail.ParseAddress(email)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user.Email = parsedEmail.Address
		user.DisplayName = displayName

		if err = userStorage.Save(); err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/user", http.StatusFound)
	}).Methods(http.MethodPost)

	r.HandleFunc("/user/password", func(w http.ResponseWriter, r *http.Request) {
		session := models.GetRequestSession(r)

		if session == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		oldPassword := r.FormValue("old_password")
		newPassword := r.FormValue("new_password")
		newPasswordConfirmation := r.FormValue("confirm_new_password")

		if len(oldPassword) == 0 || len(newPassword) == 0 || len(newPasswordConfirmation) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if newPassword != newPasswordConfirmation {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if oldPassword == newPassword {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user := userStorage.CheckUserPassword(session.User, oldPassword)
		if user == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := userStorage.SetUserPassword(user, newPassword); err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := userStorage.Save(); err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/user", http.StatusFound)
	}).Methods(http.MethodPost)
}
