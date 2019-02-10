package routes

import (
	"log"
	"net/http"
	"net/mail"
	"strconv"
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

	r.HandleFunc("/user/new", func(w http.ResponseWriter, r *http.Request) {
		session := models.GetRequestSession(r)

		if !userStorage.IsEmpty() {
			if session == nil {
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}

			if !session.User.Admin {
				w.WriteHeader(http.StatusForbidden)
				return
			}
		}

		v := view.New("layout", "newuser", config, nil)
		if err := v.Render(w); err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}).Methods(http.MethodGet)

	r.HandleFunc("/user/new", func(w http.ResponseWriter, r *http.Request) {
		session := models.GetRequestSession(r)

		if !userStorage.IsEmpty() {
			if session == nil {
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}

			if !session.User.Admin {
				w.WriteHeader(http.StatusForbidden)
				return
			}
		}

		email := r.FormValue("email")
		displayName := strings.TrimSpace(r.FormValue("display_name"))
		admin := r.FormValue("admin") == "on"
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

		user, err := models.NewUser(email, displayName, password, admin)
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = userStorage.AddUser(user)
		if err != nil {
			// User already exists.
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = userStorage.Save()
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if session == nil {
			session = models.NewSession(user)
			sessionStorage.AddSession(session)
			http.SetCookie(w, session.Cookie(false))
			http.Redirect(w, r, "/", http.StatusFound)
		} else {
			http.Redirect(w, r, "/admin", http.StatusFound)
		}
	}).Methods(http.MethodPost)

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

	r.HandleFunc("/user/delete", func(w http.ResponseWriter, r *http.Request) {
		session := models.GetRequestSession(r)
		user := session.User

		if session == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		id_t, err := strconv.ParseUint(r.FormValue("id"), 10, 32)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		id := uint32(id_t)

		if id != user.ID && !user.Admin {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		userStorage.DeleteUser(id)

		if err = userStorage.Save(); err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		sessionStorage.RemoveSessionsForUser(id)

		if user.ID == id {
			http.SetCookie(w, models.RemoveSessionCookie(false))
			http.Redirect(w, r, "/", http.StatusFound)
		} else {
			http.Redirect(w, r, "/admin", http.StatusFound)
		}
	}).Methods(http.MethodPost)
}
