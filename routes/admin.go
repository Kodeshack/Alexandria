package routes

import (
	"log"
	"net/http"
	"net/mail"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"alexandria.app/models"
	"alexandria.app/view"
)

func AdminRoutes(r *mux.Router, config *models.Config, userStorage models.UserStorage, sessionStorage *models.SessionStorage) {
	r.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		user := models.GetRequestUser(r)

		v := view.New("layout", "admin", config)
		if err := v.Render(w, user, userStorage.GetUsers()); err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}).Methods(http.MethodGet)

	r.HandleFunc("/admin/create_user", func(w http.ResponseWriter, r *http.Request) {
		session := models.GetRequestSession(r)

		email := r.FormValue("email")
		displayName := strings.TrimSpace(r.FormValue("display_name"))
		admin := r.FormValue("admin") == "on"
		password := r.FormValue("password")
		passwordConfirmation := r.FormValue("confirm_password")

		if len(email) == 0 || len(displayName) == 0 || len(password) == 0 || len(passwordConfirmation) == 0 {
			view.RenderErrorView("Email, display name, password or password confirmation empty.", http.StatusBadRequest, config, session.User, w)
			return
		}

		if password != passwordConfirmation {
			view.RenderErrorView("Passwords don't match.", http.StatusBadRequest, config, session.User, w)
			return
		}

		parsedEmail, err := mail.ParseAddress(email)
		if err != nil {
			view.RenderErrorView("Invalid email address.", http.StatusBadRequest, config, session.User, w)
			return
		}
		email = parsedEmail.Address

		user, err := models.NewUser(email, displayName, password, admin)
		if err != nil {
			log.Fatal(err)
			view.RenderErrorView("Failed to create new user.", http.StatusInternalServerError, config, session.User, w)
			return
		}

		err = userStorage.AddUser(user)
		if err != nil {
			// User already exists or, very unlikely, a UUID collision.
			view.RenderErrorView("Failed to add user to user database.", http.StatusBadRequest, config, session.User, w)
			return
		}

		err = userStorage.Save()
		if err != nil {
			log.Fatal(err)
			view.RenderErrorView("Failed to save user database.", http.StatusInternalServerError, config, session.User, w)
			return
		}

		http.Redirect(w, r, "/admin", http.StatusFound)
	}).Methods(http.MethodPost)

	r.HandleFunc("/admin/delete_user", func(w http.ResponseWriter, r *http.Request) {
		user := models.GetRequestUser(r)

		idt, err := strconv.ParseUint(r.FormValue("id"), 10, 32)
		if err != nil {
			view.RenderErrorView("Invalid user id.", http.StatusBadRequest, config, user, w)
			return
		}
		id := uint32(idt)

		userStorage.DeleteUser(id)

		if err = userStorage.Save(); err != nil {
			log.Print(err)
			view.RenderErrorView("Failed to save user database.", http.StatusInternalServerError, config, user, w)
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
