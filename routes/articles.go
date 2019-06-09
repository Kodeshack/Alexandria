package routes

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"alexandria.app/models"
	"alexandria.app/utilities"
	"alexandria.app/view"

	"github.com/gorilla/mux"
)

// ArticleRoutes sets up all HTTP routes for creating/viewing/editing routes and categories.
func ArticleRoutes(r *mux.Router, config *models.Config) {
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
		category := strings.TrimSpace(r.FormValue("category"))
		title := strings.TrimSpace(r.FormValue("title"))
		content := strings.TrimSpace(r.FormValue("content"))
		oldPath := r.FormValue("old_path")

		// For some reason when the browser POSTs data from the <textarea> it inserts `\r` before every
		// `\n` character. Because the markdown spec defines newlines as `\n` only, we need
		// to remove the offending `\r`s.
		content = strings.Replace(content, "\r", "", -1)

		article := models.NewArticle(category, title, content, config.ContentPath)

		if utilities.ArticleExists(article.Path()) && oldPath != article.Path() {
			view.RenderErrorView("Article already exists.", http.StatusBadRequest, config, user, w)
			return
		}

		err := article.Write()
		if err != nil {
			view.RenderErrorView("Failed to write article file.", http.StatusInternalServerError, config, user, w)
			return
		}

		if oldPath != article.Path() {
			err := models.RemoveArticle(oldPath)
			if err != nil {
				view.RenderErrorView("Failed to delete old article file.", http.StatusInternalServerError, config, user, w)
				return
			}
		}

		http.Redirect(w, r, "/articles/"+category+"/"+title, http.StatusFound)
	}).Methods(http.MethodPost)

	r.HandleFunc(`/articles/edit/{path:[\w\d_ /-]+}`, func(w http.ResponseWriter, r *http.Request) {
		user := models.GetRequestUser(r)

		vars := mux.Vars(r)
		path := vars["path"]
		if len(path) == 0 {
			view.RenderErrorView("", http.StatusNotFound, config, user, w)
			return
		}

		realPath := filepath.Join(config.ContentPath, path)

		stat, _ := os.Stat(realPath)

		if stat != nil && stat.IsDir() {
			view.RenderErrorView("", http.StatusNotFound, config, user, w)
			return
		}

		article, err := models.LoadArticle(config.ContentPath, filepath.Dir(path), filepath.Base(path))
		if err != nil {
			view.RenderErrorView("", http.StatusNotFound, config, user, w)
			return
		}

		v := view.New("layout", "editor", config)
		if err := v.Render(w, user, article); err != nil {
			log.Print(err)
			view.RenderErrorView("Failed to render editor view.", http.StatusInternalServerError, config, user, w)
			return
		}
	}).Methods(http.MethodGet)

	r.HandleFunc(`/articles/{path:[\w\d_ /-]+}`, func(w http.ResponseWriter, r *http.Request) {
		user := models.GetRequestUser(r)

		vars := mux.Vars(r)
		path := vars["path"]
		if len(path) == 0 {
			view.RenderErrorView("", http.StatusNotFound, config, user, w)
			return
		}

		realPath := filepath.Join(config.ContentPath, path)

		if utilities.CategoryExists(realPath) {
			category := models.NewCategory(path, realPath)
			if err := category.ScanEntries(); err != nil {
				view.RenderErrorView("Failed to read category directory.", http.StatusInternalServerError, config, user, w)
				return
			}

			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			v := view.New("layout", "category", config)
			if err := v.Render(w, user, category); err != nil {
				log.Print(err)
				view.RenderErrorView("Failed to render category view.", http.StatusInternalServerError, config, user, w)
				return
			}
		} else {
			article, err := models.LoadArticle(config.ContentPath, filepath.Dir(path), filepath.Base(path))
			if err != nil {
				view.RenderErrorView("", http.StatusNotFound, config, user, w)
				return
			}

			v := view.New("layout", "article", config)
			if err := v.Render(w, user, article); err != nil {
				log.Print(err)
				view.RenderErrorView("Failed to render article view.", http.StatusInternalServerError, config, user, w)
				return
			}
		}
	}).Methods(http.MethodGet)
}
