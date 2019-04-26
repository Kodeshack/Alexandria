package view

import (
	"log"
	"net/http"

	"alexandria.app/models"
)

type errorViewData struct {
	StatusCode int
	StatusText string
	Message    string
}

// RenderErrorView will use the provided information to render an "informative" error view.
func RenderErrorView(msg string, statusCode int, config *models.Config, user *models.User, w http.ResponseWriter) {
	data := &errorViewData{
		StatusCode: statusCode,
		StatusText: http.StatusText(statusCode),
		Message:    msg,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(statusCode)
	v := New("layout", "error", config)
	if err := v.Render(w, user, data); err != nil {
		log.Panic(err)
		return
	}
}
