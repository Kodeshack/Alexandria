package routes

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"alexandria.app/models"
	"alexandria.app/view"
)

// IndexRoutes sets up the index routes for the system.
func IndexRoutes(r *mux.Router, config *models.Config, userStorage models.UserStorage) {
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if userStorage.IsEmpty() {
			http.Redirect(w, r, "/setup", http.StatusFound)
			return
		}

		user := models.GetRequestUser(r)
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		var category *models.Category
		if category = models.NewCategory("", config.ContentPath); category != nil {
			if err := category.ScanEntries(); err != nil {
				log.Printf("Error while reading root categroy %v", err)
			}
		}

		v := view.New("layout", "index", config)

		if err := v.Render(w, user, category); err != nil {
			log.Print(err)
			view.RenderErrorView("Failed to render index layout.", http.StatusInternalServerError, config, user, w)
			return
		}
	}).Methods(http.MethodGet)

	r.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		content := []byte("User-agent: *\nDisallow: /")

		w.WriteHeader(http.StatusOK)
		w.Write(content)
	}).Methods(http.MethodGet)
}
