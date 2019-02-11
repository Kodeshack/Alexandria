package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	"alexandria.app/models"
)

func LogoutRoutes(r *mux.Router, sessionStorage *models.SessionStorage) {
	r.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		if session := models.GetRequestSession(r); session != nil {
			sessionStorage.RemoveSession(session)
			http.SetCookie(w, models.RemoveSessionCookie(false))
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}).Methods(http.MethodGet)
}
