package routes

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"alexandria.app/models"
)

func AdminRoutes(r *mux.Router, userStorage models.UserStorage) {
	r.HandleFunc("/api/user", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if r.Header.Get("content-type") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(NewErrorJSON(http.StatusBadRequest, `Content-Type must be "application/json"`))
			return
		}

		data, err := ioutil.ReadAll(r.Body)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(NewErrorJSON(http.StatusBadRequest, "Error while reading request body"))
			return
		}

		u := models.CreateUser{}
		err = json.Unmarshal(data, &u)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(NewErrorJSON(http.StatusBadRequest, "Invalid JSON"))
			return
		}

		user, err := models.NewUser(u.Email, u.DisplayName, u.Password, u.Admin)
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = userStorage.AddUser(user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(NewErrorJSON(http.StatusBadRequest, "User already exists"))
			return
		}

		err = userStorage.Save()
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		resJSON, err := user.JSON()
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(resJSON)
	})
}
