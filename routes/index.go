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
		type viewData struct {
			User     *models.User
			Category *models.Category
		}

		data := viewData{}

		if session := models.GetRequestSession(r); session != nil {
			data.User = session.User
		}

		if category := models.NewCategory("", config.ContentPath); category != nil {
			if err := category.ScanEntries(); err != nil {
				log.Printf("Error while reading root categroy %v", err)
			} else {
				data.Category = category
			}
		}

		if userStorage.IsEmpty() {
			http.Redirect(w, r, "/user/new", http.StatusFound)
			return
		}

		v := view.New("layout", "index", config, data)

		if err := v.Render(w); err != nil {
			log.Print(err)
			return
		}
	}).Methods("GET")

	r.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		content := []byte("User-agent: *\nDisallow: /")

		w.WriteHeader(http.StatusOK)
		w.Write(content)
	}).Methods("GET")
}
