package server

import (
	"log"
	"net/http"

	"alexandria.app/models"
	"alexandria.app/routes"
	"alexandria.app/view"
	"github.com/gorilla/mux"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func authedUserMiddleware(userStorage models.UserStorage) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if userStorage.IsEmpty() {
				http.Redirect(w, r, "/setup", http.StatusFound)
				return
			}

			if models.GetRequestUser(r) == nil {
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func authedAdminMiddleware(config *models.Config) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := models.GetRequestUser(r)
			if !user.Admin {
				view.RenderErrorView("", http.StatusForbidden, config, user, w)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func Start(userStorage models.UserStorage, sessionStorage *models.SessionStorage, config *models.Config) {
	r := mux.NewRouter()

	r.Use(loggingMiddleware)

	r.Use(routes.AuthMiddleWare(sessionStorage))

	authedUser := r.PathPrefix("").Subrouter()
	authedUser.Use(authedUserMiddleware(userStorage))

	authedAdmin := authedUser.PathPrefix("").Subrouter()
	authedAdmin.Use(authedAdminMiddleware(config))

	setup := r.PathPrefix("/setup").Subrouter()
	routes.SetupRoutes(setup, config, userStorage, sessionStorage)

	routes.IndexRoutes(r, config, userStorage)

	// Session-related routes.
	routes.LoginRoutes(r, config, userStorage, sessionStorage)
	routes.LogoutRoutes(authedUser, sessionStorage)

	// User-related routes.
	routes.AdminRoutes(authedAdmin, config, userStorage, sessionStorage)
	routes.UserRoutes(authedUser, config, userStorage, sessionStorage)

	// Content-related routes.
	routes.ArticleRoutes(authedUser, config)

	// Asset-related routes
	routes.AssetRoutes(r, config)

	log.Fatal(http.ListenAndServe(config.Host+config.Port, r))
}
