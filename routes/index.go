package routes

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"alexandria.app/models"
	"alexandria.app/view"
)

func IndexRoutes(r *mux.Router, config *models.Config, userStorage models.UserStorage) {
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var user *models.User

		if session := models.GetRequestSession(r); session != nil {
			user = session.User
		}

		var category *models.Category
		if category = models.NewCategory("", config.ContentPath); category != nil {
			if err := category.ScanEntries(); err != nil {
				log.Printf("Error while reading root categroy %v", err)
			}
		}

		if userStorage.IsEmpty() {
			http.Redirect(w, r, "/user/new", http.StatusFound)
			return
		}

		v := view.New("layout", "index", config)

		if err := v.Render(w, user, category); err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}).Methods(http.MethodGet)

	r.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		content := []byte("User-agent: *\nDisallow: /")

		w.WriteHeader(http.StatusOK)
		w.Write(content)
	}).Methods(http.MethodGet)
}
