package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"tobiwiki.app/models"
	"tobiwiki.app/routes"
)

const (
	ContentPrefix   = "content"
	UserStoragePath = "users.db"
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

	r := mux.NewRouter()

	r.Use(loggingMiddleware)

	routes.AdminRoutes(r, userStorage)

	routes.ArticleRoutes(r, ContentPrefix)

	log.Fatal(http.ListenAndServe(":8080", r))
}
