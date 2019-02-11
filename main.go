package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"alexandria.app/models"
	"alexandria.app/routes"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func main() {
	config := models.NewConfig()

	userStorage, err := models.LoadUserStorage(config.UserStoragePath)
	if err != nil {
		panic(err)
	}

	sessionStorage := models.NewSessionStorage()

	r := mux.NewRouter()

	r.Use(loggingMiddleware)

	r.Use(routes.AuthMiddleWare(sessionStorage))

	routes.IndexRoutes(r, config, userStorage)

	// Session-related routes.
	routes.LoginRoutes(r, config, userStorage, sessionStorage)
	routes.LogoutRoutes(r, sessionStorage)

	// User-related routes.
	routes.AdminRoutes(r, config, userStorage)
	routes.UserRoutes(r, config, userStorage, sessionStorage)

	// Content-related routes.
	routes.ArticleRoutes(r, config)

	log.Fatal(http.ListenAndServe(config.Host+config.Port, r))
}
