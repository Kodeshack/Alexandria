package routes

import (
	"net/http"

	"alexandria.app/models"
	"github.com/gorilla/mux"
)

func AssetRoutes(r *mux.Router, config *models.Config) {
	r.PathPrefix(`/assets/{f:[\-\d\w]+\.[js|css]}`).Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir(config.AssetPath))))
}
