package routes

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"

	"tobiwiki.app/models"
)

func AdminRoutes(r *mux.Router) {
	r.HandleFunc("/api/user", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if r.Header.Get("content-type") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		data, err := ioutil.ReadAll(r.Body)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		u := models.CreateUser{}
		err = json.Unmarshal(data, &u)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user, err := models.NewUser(u.Email, u.DisplayName, u.Password, u.Admin)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		resJSON, err := user.JSON()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(resJSON)
	})
}
