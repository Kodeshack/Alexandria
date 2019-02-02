package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"alexandria.app/models"
	"alexandria.app/routes"
)

const (
	ContentPrefix   = "content"
	UserStoragePath = "users.db"
	TemplateDir     = "view/templates"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func main() {
	userStorage, err := models.LoadUserStorage(UserStoragePath)
	if err != nil {
		panic(err)
	}

	sessionStorage := models.NewSessionStorage()

	r := mux.NewRouter()

	r.Use(loggingMiddleware)

	r.Use(routes.AuthMiddleWare(sessionStorage))

	routes.AdminRoutes(r, TemplateDir, userStorage, sessionStorage)

	routes.ArticleRoutes(r, ContentPrefix, TemplateDir)

	routes.IndexRoutes(r, TemplateDir)

	log.Fatal(http.ListenAndServe(":8080", r))
}
