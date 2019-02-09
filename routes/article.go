package routes

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"alexandria.app/models"
	"alexandria.app/view"

	"github.com/gorilla/mux"
)

func ArticleRoutes(r *mux.Router, config *models.Config) {
	r.HandleFunc("/article/new", func(w http.ResponseWriter, r *http.Request) {
		if models.GetRequestSession(r) == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		v := view.New("layout", "editor", config, nil)

		if err := v.Render(w); err != nil {
			log.Print(err)
			return
		}
	})

	r.HandleFunc("/article/save", func(w http.ResponseWriter, r *http.Request) {
		if models.GetRequestSession(r) == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		title := r.FormValue("title")

		// For some reason when the browser POSTs data from the <textarea> it inserts `\r` before every
		// `\n` character. Because the markdown spec defines newlines as `\n` only, we need
		// to remove the offending `\r`s.
		content := strings.Replace(r.FormValue("content"), "\r", "", -1)

		dir := filepath.Join(config.ContentPath, filepath.Dir(title))
		fileName := filepath.Base(title)

		article := models.NewArticle(fileName, content, dir)

		err := article.Write()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/"+title, http.StatusFound)
	})

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

		realPath := filepath.Join(config.ContentPath, path)

		stat, _ := os.Stat(realPath)

		if stat != nil && stat.IsDir() {
			category := models.NewCategory(path, realPath)
			if err := category.ScanEntries(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			v := view.New("layout", "category", config, category)
			if err := v.Render(w); err != nil {
				log.Print(err)
				return
			}
		} else {
			article, err := models.LoadArticle(realPath + ".md")
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			data, err := article.ContentHTML()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			w.Write(data)
		}
	})
}
