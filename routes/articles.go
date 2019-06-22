package routes

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"alexandria.app/articledb"
	"alexandria.app/models"
	"alexandria.app/view"

	"github.com/gorilla/mux"
)

// ArticleRoutes sets up all HTTP routes for creating/viewing/editing routes and categories.
func ArticleRoutes(r *mux.Router, articledb articledb.ArticleDB, config *models.Config) {
	r.HandleFunc("/articles/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusFound)
	}).Methods(http.MethodGet)

	r.HandleFunc("/articles/new", func(w http.ResponseWriter, r *http.Request) {
		user := models.GetRequestUser(r)

		v := view.New("layout", "editor", config)
		if err := v.Render(w, user, nil); err != nil {
			log.Print(err)
			view.RenderErrorView("Failed to render editor view.", http.StatusInternalServerError, config, user, w)
			return
		}
	}).Methods(http.MethodGet)

	r.HandleFunc("/articles/save", func(w http.ResponseWriter, r *http.Request) {
		user := models.GetRequestUser(r)
		title := strings.TrimSpace(r.FormValue("title"))
		content := strings.TrimSpace(r.FormValue("content"))

		// For some reason when the browser POSTs data from the <textarea> it inserts `\r` before every
		// `\n` character. Because the markdown spec defines newlines as `\n` only, we need
		// to remove the offending `\r`s.
		content = strings.Replace(content, "\r", "", -1)

		dir := filepath.Dir(title)
		fileName := filepath.Base(title)

		article := models.NewArticle(fileName, content, dir)

		err := articledb.Write(article)
		if err != nil {
			view.RenderErrorView("Failed to write article file.", http.StatusInternalServerError, config, user, w)
			return
		}

		http.Redirect(w, r, "/articles/"+title, http.StatusFound)
	}).Methods(http.MethodPost)

	r.HandleFunc(`/articles/{path:[\w\d_ /-]+}`, func(w http.ResponseWriter, r *http.Request) {
		user := models.GetRequestUser(r)

		vars := mux.Vars(r)
		path := vars["path"]
		if len(path) == 0 {
			view.RenderErrorView("", http.StatusNotFound, config, user, w)
			return
		}

		article, category, err := articledb.Load(path)
		if err != nil {
			view.RenderErrorView("Failed to read file/dir.", http.StatusInternalServerError, config, user, w)
			return
		}

		if article == nil && category == nil {
			view.RenderErrorView("", http.StatusNotFound, config, user, w)
			return
		}

		if category != nil {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			v := view.New("layout", "category", config)
			if err := v.Render(w, user, category); err != nil {
				log.Print(err)
				view.RenderErrorView("Failed to render category view.", http.StatusInternalServerError, config, user, w)
				return
			}
		} else {
			body, err := article.ContentHTML()
			if err != nil {
				view.RenderErrorView("Failed to render content as HTML.", http.StatusInternalServerError, config, user, w)
				return
			}

			v := view.New("layout", "article", config)
			if err := v.Render(w, user, template.HTML(body)); err != nil {
				log.Print(err)
				view.RenderErrorView("Failed to render article view.", http.StatusInternalServerError, config, user, w)
				return
			}
		}
	}).Methods(http.MethodGet)
}
