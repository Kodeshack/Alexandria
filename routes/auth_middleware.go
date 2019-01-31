package routes

import (
	"context"
	"net/http"

	"alexandria.app/models"
	"github.com/gorilla/mux"
)

func AuthMiddleWare(sessionStorage *models.SessionStorage) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessionCookie, err := r.Cookie(models.SessionCookieName)

			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			session := sessionStorage.GetSession(sessionCookie.Value)

			ctx := context.WithValue(r.Context(), models.SessionContextKey, session)

			req := r.WithContext(ctx)

			next.ServeHTTP(w, req)
		})
	}
}
