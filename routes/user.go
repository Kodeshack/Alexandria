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

// UserRoutes sets up all HTTP routes required for user management.
func UserRoutes(r *mux.Router, config *models.Config, userStorage models.UserStorage, sessionStorage *models.SessionStorage) {
	r.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		user := models.GetRequestUser(r)

		v := view.New("layout", "user", config)
		if err := v.Render(w, user, user); err != nil {
			log.Print(err)
			view.RenderErrorView("Failed to render user view.", http.StatusInternalServerError, config, user, w)
			return
		}
	}).Methods(http.MethodGet)

	r.HandleFunc("/user/new", func(w http.ResponseWriter, r *http.Request) {
		user := models.GetRequestUser(r)

		if !user.Admin {
			view.RenderErrorView("", http.StatusForbidden, config, user, w)
			return
		}

		v := view.New("layout", "newuser", config)
		if err := v.Render(w, user, nil); err != nil {
			log.Print(err)
			view.RenderErrorView("Failed to render new user view.", http.StatusInternalServerError, config, user, w)
			return
		}
	}).Methods(http.MethodGet)

	r.HandleFunc("/user/update", func(w http.ResponseWriter, r *http.Request) {
		user := models.GetRequestUser(r)
		email := r.FormValue("email")
		displayName := strings.TrimSpace(r.FormValue("display_name"))

		if len(email) == 0 || len(displayName) == 0 {
			view.RenderErrorView("Email or display name empty.", http.StatusBadRequest, config, user, w)
			return
		}

		parsedEmail, err := mail.ParseAddress(email)
		if err != nil {
			view.RenderErrorView("Invalid email address", http.StatusBadRequest, config, user, w)
			return
		}

		user.Email = parsedEmail.Address
		user.DisplayName = displayName

		if err = userStorage.Save(); err != nil {
			log.Print(err)
			view.RenderErrorView("Failed to save user to user database.", http.StatusInternalServerError, config, user, w)
			return
		}

		http.Redirect(w, r, "/user", http.StatusFound)
	}).Methods(http.MethodPost)

	r.HandleFunc("/user/password", func(w http.ResponseWriter, r *http.Request) {
		session := models.GetRequestSession(r)

		oldPassword := r.FormValue("old_password")
		newPassword := r.FormValue("new_password")
		newPasswordConfirmation := r.FormValue("confirm_new_password")

		if len(oldPassword) == 0 || len(newPassword) == 0 || len(newPasswordConfirmation) == 0 {
			view.RenderErrorView("Old password, new password or new password confirmation empty.", http.StatusBadRequest, config, session.User, w)
			return
		}

		if newPassword != newPasswordConfirmation {
			view.RenderErrorView("New passwords don't match.", http.StatusBadRequest, config, session.User, w)
			return
		}

		if oldPassword == newPassword {
			view.RenderErrorView("New password is the same as old password.", http.StatusBadRequest, config, session.User, w)
			return
		}

		user := userStorage.CheckUserPassword(session.User, oldPassword)
		if user == nil {
			view.RenderErrorView("Incorrect old password.", http.StatusBadRequest, config, session.User, w)
			return
		}

		if err := userStorage.SetUserPassword(user, newPassword); err != nil {
			log.Print(err)
			view.RenderErrorView("Failed to update password.", http.StatusInternalServerError, config, session.User, w)
			return
		}

		if err := userStorage.Save(); err != nil {
			log.Print(err)
			view.RenderErrorView("Failed to save user database.", http.StatusInternalServerError, config, session.User, w)
			return
		}

		http.Redirect(w, r, "/user", http.StatusFound)
	}).Methods(http.MethodPost)

	r.HandleFunc("/user/delete", func(w http.ResponseWriter, r *http.Request) {
		user := models.GetRequestUser(r)

		idt, err := strconv.ParseUint(r.FormValue("id"), 10, 32)
		if err != nil {
			view.RenderErrorView("Invalid user id.", http.StatusBadRequest, config, user, w)
			return
		}
		id := uint32(idt)

		if id != user.ID {
			view.RenderErrorView("Can't delete currently logged in user.", http.StatusForbidden, config, user, w)
			return
		}

		userStorage.DeleteUser(id)

		if err = userStorage.Save(); err != nil {
			log.Print(err)
			view.RenderErrorView("Failed to save user database.", http.StatusInternalServerError, config, user, w)
			return
		}

		sessionStorage.RemoveSessionsForUser(id)

		http.SetCookie(w, models.RemoveSessionCookie(false))
		http.Redirect(w, r, "/", http.StatusFound)
	}).Methods(http.MethodPost)
}
