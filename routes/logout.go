package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	"alexandria.app/models"
)

// LogoutRoutes sets up all HTTP routes for user logout.
func LogoutRoutes(r *mux.Router, sessionStorage *models.SessionStorage) {
	r.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		session := models.GetRequestSession(r)
		sessionStorage.RemoveSession(session)
		http.SetCookie(w, models.RemoveSessionCookie(false))
		http.Redirect(w, r, "/", http.StatusFound)
	}).Methods(http.MethodGet)
}
