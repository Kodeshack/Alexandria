package routes

import (
	"html/template"
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
	r.HandleFunc("/articles/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusFound)
	}).Methods(http.MethodGet)

	r.HandleFunc("/articles/new", func(w http.ResponseWriter, r *http.Request) {
		user := models.GetRequestUser(r)
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		v := view.New("layout", "editor", config)
		if err := v.Render(w, user, nil); err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}).Methods(http.MethodGet)

	r.HandleFunc("/articles/save", func(w http.ResponseWriter, r *http.Request) {
		if models.GetRequestSession(r) == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		title := strings.TrimSpace(r.FormValue("title"))
		content := strings.TrimSpace(r.FormValue("content"))

		// For some reason when the browser POSTs data from the <textarea> it inserts `\r` before every
		// `\n` character. Because the markdown spec defines newlines as `\n` only, we need
		// to remove the offending `\r`s.
		content = strings.Replace(content, "\r", "", -1)

		dir := filepath.Join(config.ContentPath, filepath.Dir(title))
		fileName := filepath.Base(title)

		article := models.NewArticle(fileName, content, dir)

		err := article.Write()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/articles/"+title, http.StatusFound)
	}).Methods(http.MethodPost)

	r.HandleFunc(`/articles/{path:[\w\d_/-]+}`, func(w http.ResponseWriter, r *http.Request) {
		user := models.GetRequestUser(r)
		if user == nil {
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
			v := view.New("layout", "category", config)
			if err := v.Render(w, user, category); err != nil {
				log.Print(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			article, err := models.LoadArticle(realPath + ".md")
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			body, err := article.ContentHTML()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			v := view.New("layout", "article", config)
			if err := v.Render(w, user, template.HTML(body)); err != nil {
				log.Print(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	}).Methods(http.MethodGet)
}
