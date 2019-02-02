package routes

import (
	"net/http"
	"path/filepath"

	"alexandria.app/models"

	"github.com/gorilla/mux"
)

func ArticleRoutes(r *mux.Router, contentPath string) {
	r.HandleFunc("/{path:[\\w\\d_/-]+}", func(w http.ResponseWriter, r *http.Request) {
		if models.GetRequestSession(r) == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		vars := mux.Vars(r)
		path := vars["path"]
		if len(path) == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		article, err := models.LoadArticle(filepath.Join(contentPath, path+".md"))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		data, err := article.ContentHTML()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})
}
